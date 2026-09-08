package cli

import (
	"fmt"
	"io"
)

const defaultProgressInterval = 1000

func reportProgress(writer io.Writer, phase string, completed, total, progressInterval int) {
	if completed != 0 && completed != total && completed%progressInterval != 0 {
		return
	}

	percentage := 100.0
	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}
	fmt.Fprintf(writer, "%s progress: %d/%d (%.1f%%)\n", phase, completed, total, percentage)
}
