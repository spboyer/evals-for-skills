package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// SessionFile represents a session log file on disk.
type SessionFile struct {
	Path      string
	Name      string
	Size      int64
	ModTime   time.Time
	NumEvents int
}

// ListSessions finds .jsonl session log files in dir.
func ListSessions(dir string) ([]SessionFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading session directory: %w", err)
	}

	var files []SessionFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), "-session.jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}

		path := filepath.Join(dir, e.Name())
		n, _ := countLines(path)
		files = append(files, SessionFile{
			Path:      path,
			Name:      e.Name(),
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			NumEvents: n,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})

	return files, nil
}

func countLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()
	n := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n++
	}
	return n, scanner.Err()
}

// ReadEvents parses all events from a session log file.
func ReadEvents(path string) ([]Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening session file: %w", err)
	}
	defer func() { _ = f.Close() }()

	var events []Event
	scanner := bufio.NewScanner(f)
	// Increase buffer for large lines.
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		var ev Event
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue // skip malformed lines
		}
		events = append(events, ev)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading session file: %w", err)
	}
	return events, nil
}

// RenderTimeline writes a human-readable session timeline to w.
func RenderTimeline(w io.Writer, events []Event) {
	if len(events) == 0 {
		_, _ = fmt.Fprintln(w, "No events found.")
		return
	}

	_, _ = fmt.Fprintln(w, "═══════════════════════════════════════════════════════")
	_, _ = fmt.Fprintln(w, " SESSION TIMELINE")
	_, _ = fmt.Fprintln(w, "═══════════════════════════════════════════════════════")
	_, _ = fmt.Fprintln(w)

	start := events[0].Timestamp
	for _, ev := range events {
		elapsed := ev.Timestamp.Sub(start)
		ts := formatDuration(elapsed)

		switch ev.Type {
		case EventSessionStart:
			model, _ := ev.Data["model"].(string) //nolint:errcheck // zero-value on failure is fine
			engine, _ := ev.Data["engine"].(string) //nolint:errcheck // zero-value on failure is fine
			taskCount := jsonNumber(ev.Data["task_count"])
			_, _ = fmt.Fprintf(w, "[%s] 🚀 Session started  model=%s  engine=%s  tasks=%d\n", ts, model, engine, taskCount)

		case EventTaskStart:
			name, _ := ev.Data["task_name"].(string) //nolint:errcheck // zero-value on failure is fine
			num := jsonNumber(ev.Data["task_num"])
			total := jsonNumber(ev.Data["total_tasks"])
			_, _ = fmt.Fprintf(w, "[%s] ▶  Task %d/%d: %s\n", ts, num, total, name)

		case EventGraderResult:
			grader, _ := ev.Data["grader_name"].(string) //nolint:errcheck // zero-value on failure is fine
			passed, _ := ev.Data["passed"].(bool) //nolint:errcheck // zero-value on failure is fine
			score := jsonFloat(ev.Data["score"])
			icon := "✗"
			if passed {
				icon = "✓"
			}
			_, _ = fmt.Fprintf(w, "[%s]    %s Grader %s  score=%.2f\n", ts, icon, grader, score)

		case EventTaskComplete:
			name, _ := ev.Data["task_name"].(string) //nolint:errcheck // zero-value on failure is fine
			status, _ := ev.Data["status"].(string) //nolint:errcheck // zero-value on failure is fine
			dur := jsonNumber(ev.Data["duration_ms"])
			icon := "✓"
			if status != "passed" {
				icon = "✗"
			}
			_, _ = fmt.Fprintf(w, "[%s] %s  Task complete: %s [%s] (%dms)\n", ts, icon, name, status, dur)

		case EventError:
			msg, _ := ev.Data["message"].(string) //nolint:errcheck // zero-value on failure is fine
			_, _ = fmt.Fprintf(w, "[%s] ❌ Error: %s\n", ts, msg)

		case EventSessionEnd:
			total := jsonNumber(ev.Data["total_tests"])
			passed := jsonNumber(ev.Data["passed"])
			failed := jsonNumber(ev.Data["failed"])
			dur := jsonNumber(ev.Data["duration_ms"])
			_, _ = fmt.Fprintf(w, "[%s] 🏁 Session complete  %d/%d passed  %d failed  (%dms)\n",
				ts, passed, total, failed, dur)

		default:
			_, _ = fmt.Fprintf(w, "[%s] %s %v\n", ts, ev.Type, ev.Data)
		}
	}
	_, _ = fmt.Fprintln(w)
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%6dms", d.Milliseconds())
	}
	return fmt.Sprintf("%6.1fs", d.Seconds())
}

// jsonNumber extracts a number from a JSON-decoded interface{} (float64 or json.Number).
func jsonNumber(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case json.Number:
		i, _ := n.Int64() //nolint:errcheck // fallback to zero is acceptable
		return int(i)
	}
	return 0
}

func jsonFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int:
		return float64(n)
	case json.Number:
		f, _ := n.Float64() //nolint:errcheck // fallback to zero is acceptable
		return f
	}
	return 0
}
