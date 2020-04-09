package textformat

import (
	"context"
	"fmt"
	"io"
)

// AttachVolume writes a success message in the writer
func (d *Displayer) AttachVolume(ctx context.Context, w io.Writer) error {
	_, err := fmt.Fprintln(w, "volume attached")
	return err
}
