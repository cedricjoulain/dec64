package dec64

import (
	"encoding/binary"
	"io"
)

// ListFromReader returns list of dec64 from reader
func ListFromReader(r io.Reader) (values []Dec64, err error) {
	// reading buffer
	// size as a Dec64 aka int64
	buff := make([]byte, 8)
	// small capacity to start
	values = make([]Dec64, 0, 16)
	for {
		// Read one more time
		_, err = io.ReadAtLeast(r, buff, 8)
		if err != nil {
			return
		}
		// unsafe faster ?
		values = append(values, (Dec64)(binary.LittleEndian.Uint64(buff)))
	}
}

// ListToWriter sends list of dec64 to writer
func ListToWriter(w io.Writer, values []Dec64) (err error) {
	// writing buffer
	// size as a Dec64 aka int64
	buff := make([]byte, 8)
	for _, v := range values {
		binary.LittleEndian.PutUint64(buff, uint64(v))
		_, err = w.Write(buff)
		if err != nil {
			return
		}
	}
	return
}
