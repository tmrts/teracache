package payload

import (
	"bytes"
	"io"
)

type Key string

type Payload []byte

func (p *Payload) Reader() io.Reader {
	// TODO(tmrts): cache entries are immutable, stop copying byte buffers around.
	return bytes.NewReader(*p)
}
