package dev

import (
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/spboyer/waza/internal/skill"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// LinkIssue describes a single link problem found during validation.
type LinkIssue struct {
	Source string // source file (relative to skill dir)
	Target string // link target
	Reason string // human-readable description
}

// LinkResult holds the output from link validation checks.
type LinkResult struct {
	BrokenLinks    []LinkIssue
	DirectoryLinks []LinkIssue
	ScopeEscapes   []LinkIssue
	DeadURLs       []LinkIssue
	OrphanedFiles  []string
	TotalLinks     int
	ValidLinks     int
	Issues         []Issue
}

// Passed returns true when no link errors were found.
func (r *LinkResult) Passed() bool {
	return len(r.BrokenLinks) == 0 &&
		len(r.DirectoryLinks) == 0 &&
		len(r.ScopeEscapes) == 0 &&
		len(r.DeadURLs) == 0 &&
		len(r.OrphanedFiles) == 0
}

// LinkScorer validates links in a skill's markdown files.
type LinkScorer struct{}

// linkHTTPClient is the shared HTTP client for external URL checks.
var linkHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("too many redirects")
		}
		return nil
	},
}

type extractedLink struct {
	target string
}

type fileEntry struct {
	relPath string
	links   []extractedLink
}

// Score runs link validation on the skill directory.
func (LinkScorer) Score(sk *skill.Skill) *LinkResult {
	if sk == nil || sk.Path == "" {
		return nil
	}
	skillDir := filepath.Dir(sk.Path)
	r := &LinkResult{}

	mdFiles := collectMarkdownFiles(skillDir)
	if len(mdFiles) == 0 {
		return r
	}

	// Extract links from each file.
	var entries []fileEntry
	uniqueExtURLs := make(map[string][]string) // deduped URL → source files
	for _, f := range mdFiles {
		relPath, _ := filepath.Rel(skillDir, f)
		relPath = filepath.ToSlash(relPath)
		links := extractLinksFromFile(f)
		entries = append(entries, fileEntry{relPath: relPath, links: links})

		for _, l := range links {
			if isExternalURL(l.target) {
				clean := stripFragment(l.target)
				uniqueExtURLs[clean] = appendUnique(uniqueExtURLs[clean], relPath)
			}
		}
	}

	// Validate local links.
	for _, fe := range entries {
		for _, l := range fe.links {
			target := l.target
			if shouldSkipLink(target) || isExternalURL(target) {
				continue
			}

			localTarget := stripFragment(target)
			if localTarget == "" {
				continue // fragment-only
			}

			r.TotalLinks++

			sourceDir := filepath.Dir(filepath.Join(skillDir, filepath.FromSlash(fe.relPath)))
			resolved := filepath.Clean(filepath.Join(sourceDir, filepath.FromSlash(localTarget)))

			if !isWithinDir(resolved, skillDir) {
				r.ScopeEscapes = append(r.ScopeEscapes, LinkIssue{
					Source: fe.relPath, Target: target, Reason: "link escapes skill directory",
				})
				continue
			}

			info, err := os.Stat(resolved)
			if err != nil {
				r.BrokenLinks = append(r.BrokenLinks, LinkIssue{
					Source: fe.relPath, Target: target, Reason: "target does not exist",
				})
				continue
			}

			if info.IsDir() {
				r.DirectoryLinks = append(r.DirectoryLinks, LinkIssue{
					Source: fe.relPath, Target: target, Reason: "target is a directory, not a file",
				})
				continue
			}
		}
	}

	// Validate external URLs concurrently (one check per unique URL).
	r.TotalLinks += len(uniqueExtURLs)
	r.DeadURLs = checkExternalURLsDedup(uniqueExtURLs)

	// Compute valid links.
	problems := len(r.BrokenLinks) + len(r.DirectoryLinks) + len(r.ScopeEscapes) + len(r.DeadURLs)
	r.ValidLinks = r.TotalLinks - problems
	if r.ValidLinks < 0 {
		r.ValidLinks = 0
	}

	// Orphaned file detection.
	r.OrphanedFiles = findOrphanedFiles(skillDir, entries)

	buildLinkSummaryIssues(r)
	return r
}

// collectMarkdownFiles walks dir and returns paths to .md and .mdx files.
func collectMarkdownFiles(dir string) []string {
	var files []string
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".md" || ext == ".mdx" {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// extractLinksFromFile parses a markdown file and returns all link targets.
func extractLinksFromFile(filePath string) []extractedLink {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return extractLinksFromSource(data)
}

// extractLinksFromSource parses markdown bytes and extracts link/image destinations.
func extractLinksFromSource(source []byte) []extractedLink {
	md := goldmark.New()
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	var links []extractedLink
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Link:
			links = append(links, extractedLink{target: string(v.Destination)})
		case *ast.Image:
			links = append(links, extractedLink{target: string(v.Destination)})
		case *ast.AutoLink:
			target := string(v.Label(source))
			if len(v.Protocol) > 0 && !strings.HasPrefix(target, string(v.Protocol)) {
				target = string(v.Protocol) + target
			}
			links = append(links, extractedLink{target: target})
		}
		return ast.WalkContinue, nil
	})
	return links
}

// isExternalURL returns true for http:// and https:// URLs.
func isExternalURL(target string) bool {
	return strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://")
}

// shouldSkipLink returns true for link schemes that should not be validated.
func shouldSkipLink(target string) bool {
	return strings.HasPrefix(target, "mailto:") || strings.HasPrefix(target, "mdc:")
}

// stripFragment removes the #fragment portion of a URL or path.
func stripFragment(target string) string {
	if idx := strings.Index(target, "#"); idx >= 0 {
		return target[:idx]
	}
	return target
}

// isWithinDir returns true if path is inside dir (or is dir itself).
func isWithinDir(path, dir string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// checkExternalURLsDedup validates unique external URLs with a goroutine pool of 5.
// Returns one LinkIssue per dead URL (using the first source file for reporting).
func checkExternalURLsDedup(urls map[string][]string) []LinkIssue {
	if len(urls) == 0 {
		return nil
	}

	type deadResult struct {
		url    string
		source string
		reason string
	}

	var mu sync.Mutex
	var dead []deadResult
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for u, sources := range urls {
		src := sources[0]
		wg.Add(1)
		go func(u, src string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if isDead, reason := checkSingleURL(u); isDead {
				mu.Lock()
				dead = append(dead, deadResult{url: u, source: src, reason: reason})
				mu.Unlock()
			}
		}(u, src)
	}

	wg.Wait()

	issues := make([]LinkIssue, len(dead))
	for i, d := range dead {
		issues[i] = LinkIssue{Source: d.source, Target: d.url, Reason: d.reason}
	}
	return issues
}

// skipSSRFCheck disables private-IP rejection for testing with httptest servers.
var skipSSRFCheck bool

// isPrivateIP returns true if the IP is loopback, private, or link-local.
func isPrivateIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

// checkSingleURL tries HTTP HEAD then falls back to GET.
func checkSingleURL(rawURL string) (dead bool, reason string) {
	if !skipSSRFCheck {
		// SSRF protection: reject URLs targeting private/loopback/link-local IPs.
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return true, err.Error()
		}
		host := parsed.Hostname()
		ips, err := net.LookupHost(host)
		if err != nil {
			return true, fmt.Sprintf("DNS lookup failed: %v", err)
		}
		for _, ipStr := range ips {
			if ip := net.ParseIP(ipStr); ip != nil && isPrivateIP(ip) {
				return true, fmt.Sprintf("URL resolves to private/loopback address (%s)", ipStr)
			}
		}
	}

	req, err := http.NewRequest(http.MethodHead, rawURL, nil)
	if err != nil {
		return true, err.Error()
	}
	req.Header.Set("User-Agent", "waza-link-checker/1.0")

	resp, err := linkHTTPClient.Do(req)
	if err == nil {
		_ = resp.Body.Close()
		if resp.StatusCode < 400 {
			return false, ""
		}
		// Some servers reject HEAD; fall back to GET.
		if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusForbidden {
			return checkSingleURLGet(rawURL)
		}
		return true, fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return checkSingleURLGet(rawURL)
}

func checkSingleURLGet(rawURL string) (bool, string) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return true, err.Error()
	}
	req.Header.Set("User-Agent", "waza-link-checker/1.0")

	resp, err := linkHTTPClient.Do(req)
	if err != nil {
		return true, err.Error()
	}
	_ = resp.Body.Close()
	if resp.StatusCode >= 400 {
		return true, fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return false, ""
}

// findOrphanedFiles does BFS from SKILL.md following local links transitively,
// then reports files in references/ that are not reachable.
func findOrphanedFiles(skillDir string, entries []fileEntry) []string {
	refsDir := filepath.Join(skillDir, "references")
	info, err := os.Stat(refsDir)
	if err != nil || !info.IsDir() {
		return nil
	}

	var refFiles []string
	_ = filepath.WalkDir(refsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(skillDir, path)
		refFiles = append(refFiles, filepath.ToSlash(rel))
		return nil
	})
	if len(refFiles) == 0 {
		return nil
	}

	// Build normalized link map for BFS.
	normLinkMap := make(map[string][]string)
	for _, fe := range entries {
		key := normalizePath(fe.relPath)
		var targets []string
		for _, l := range fe.links {
			if shouldSkipLink(l.target) || isExternalURL(l.target) {
				continue
			}
			localTarget := stripFragment(l.target)
			if localTarget == "" {
				continue
			}
			sourceDir := filepath.Dir(filepath.Join(skillDir, filepath.FromSlash(fe.relPath)))
			resolved := filepath.Clean(filepath.Join(sourceDir, filepath.FromSlash(localTarget)))
			rel, err := filepath.Rel(skillDir, resolved)
			if err != nil {
				continue
			}
			targets = append(targets, normalizePath(filepath.ToSlash(rel)))
		}
		normLinkMap[key] = targets
	}

	// BFS from SKILL.md.
	reachable := make(map[string]bool)
	visited := make(map[string]bool)
	queue := []string{normalizePath("SKILL.md")}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true
		reachable[current] = true
		for _, t := range normLinkMap[current] {
			if !visited[t] {
				queue = append(queue, t)
			}
		}
	}

	var orphaned []string
	for _, rf := range refFiles {
		if !reachable[normalizePath(rf)] {
			orphaned = append(orphaned, rf)
		}
	}
	return orphaned
}

// normalizePath normalizes a path for comparison (case-insensitive on Windows).
func normalizePath(p string) string {
	p = filepath.ToSlash(p)
	if runtime.GOOS == "windows" {
		p = strings.ToLower(p)
	}
	return p
}

// appendUnique appends item to slice only if not already present.
func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

// buildLinkSummaryIssues populates the Issues field from categorized problems.
func buildLinkSummaryIssues(r *LinkResult) {
	if len(r.BrokenLinks) > 0 {
		r.Issues = append(r.Issues, Issue{
			Rule: "link-broken", Message: fmt.Sprintf("%d broken link(s) found", len(r.BrokenLinks)), Severity: "error",
		})
	}
	if len(r.DirectoryLinks) > 0 {
		r.Issues = append(r.Issues, Issue{
			Rule: "link-directory", Message: fmt.Sprintf("%d link(s) point to directories instead of files", len(r.DirectoryLinks)), Severity: "warning",
		})
	}
	if len(r.ScopeEscapes) > 0 {
		r.Issues = append(r.Issues, Issue{
			Rule: "link-scope", Message: fmt.Sprintf("%d link(s) escape the skill directory", len(r.ScopeEscapes)), Severity: "error",
		})
	}
	if len(r.DeadURLs) > 0 {
		r.Issues = append(r.Issues, Issue{
			Rule: "link-dead-url", Message: fmt.Sprintf("%d dead external URL(s) found", len(r.DeadURLs)), Severity: "warning",
		})
	}
	if len(r.OrphanedFiles) > 0 {
		r.Issues = append(r.Issues, Issue{
			Rule: "link-orphan", Message: fmt.Sprintf("%d file(s) in references/ not linked from SKILL.md", len(r.OrphanedFiles)), Severity: "warning",
		})
	}
}
