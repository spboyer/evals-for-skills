package spinner

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

var frames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Start displays an animated spinner with the given message on w.
// Call the returned function to stop the spinner and clear the line.
func Start(w io.Writer, message string) (stop func()) {
	done := make(chan struct{})
	cleared := make(chan struct{})
	var stopOnce sync.Once
	go func() {
		i := 0
		for {
			select {
			case <-done:
				fmt.Fprintf(w, "\r%s\r", strings.Repeat(" ", len(message)+2)) //nolint:errcheck
				close(cleared)
				return
			case <-time.After(80 * time.Millisecond):
				fmt.Fprintf(w, "\r%s %s", frames[i%len(frames)], message) //nolint:errcheck
				i++
			}
		}
	}()
	return func() {
		stopOnce.Do(func() {
			close(done)
		})
		<-cleared
	}
}
