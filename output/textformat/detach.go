package textformat

import (
	"context"
	"fmt"
	"io"
)

// DetachVolume writes a success message to the writer
func (d *Displayer) DetachVolume(ctx context.Context, w io.Writer) error {
	_, err := fmt.Fprintln(w, "volume detached")
	return err
}
