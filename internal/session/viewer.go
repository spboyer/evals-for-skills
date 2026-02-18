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
		n, _ := countLines(path) //nolint:errcheck
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
	defer f.Close() //nolint:errcheck
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
	defer f.Close() //nolint:errcheck

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
//
//nolint:errcheck // display-only writes; errors are not actionable
func RenderTimeline(w io.Writer, events []Event) {
	if len(events) == 0 {
		fmt.Fprintln(w, "No events found.")
		return
	}

	fmt.Fprintln(w, "═══════════════════════════════════════════════════════")
	fmt.Fprintln(w, " SESSION TIMELINE")
	fmt.Fprintln(w, "═══════════════════════════════════════════════════════")
	fmt.Fprintln(w)

	start := events[0].Timestamp
	for _, ev := range events {
		elapsed := ev.Timestamp.Sub(start)
		ts := formatDuration(elapsed)

		switch ev.Type {
		case EventSessionStart:
			model, _ := ev.Data["model"].(string)   //nolint:errcheck
			engine, _ := ev.Data["engine"].(string) //nolint:errcheck
			taskCount := jsonNumber(ev.Data["task_count"])
			fmt.Fprintf(w, "[%s] 🚀 Session started  model=%s  engine=%s  tasks=%d\n", ts, model, engine, taskCount)

		case EventTaskStart:
			name, _ := ev.Data["task_name"].(string) //nolint:errcheck
			num := jsonNumber(ev.Data["task_num"])
			total := jsonNumber(ev.Data["total_tasks"])
			fmt.Fprintf(w, "[%s] ▶  Task %d/%d: %s\n", ts, num, total, name)

		case EventGraderResult:
			grader, _ := ev.Data["grader_name"].(string) //nolint:errcheck
			passed, _ := ev.Data["passed"].(bool)        //nolint:errcheck
			score := jsonFloat(ev.Data["score"])
			icon := "✗"
			if passed {
				icon = "✓"
			}
			fmt.Fprintf(w, "[%s]    %s Grader %s  score=%.2f\n", ts, icon, grader, score)

		case EventTaskComplete:
			name, _ := ev.Data["task_name"].(string) //nolint:errcheck
			status, _ := ev.Data["status"].(string)  //nolint:errcheck
			dur := jsonNumber(ev.Data["duration_ms"])
			icon := "✓"
			if status != "passed" {
				icon = "✗"
			}
			fmt.Fprintf(w, "[%s] %s  Task complete: %s [%s] (%dms)\n", ts, icon, name, status, dur)

		case EventError:
			msg, _ := ev.Data["message"].(string) //nolint:errcheck
			fmt.Fprintf(w, "[%s] ❌ Error: %s\n", ts, msg)

		case EventSessionEnd:
			total := jsonNumber(ev.Data["total_tests"])
			passed := jsonNumber(ev.Data["passed"])
			failed := jsonNumber(ev.Data["failed"])
			dur := jsonNumber(ev.Data["duration_ms"])
			fmt.Fprintf(w, "[%s] 🏁 Session complete  %d/%d passed  %d failed  (%dms)\n",
				ts, passed, total, failed, dur)

		default:
			fmt.Fprintf(w, "[%s] %s %v\n", ts, ev.Type, ev.Data)
		}
	}
	fmt.Fprintln(w)
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
		i, _ := n.Int64() //nolint:errcheck
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
		f, _ := n.Float64() //nolint:errcheck
		return f
	}
	return 0
}
