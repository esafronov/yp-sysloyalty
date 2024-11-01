package helpers

import (
	"bytes"
	"compress/gzip"
)

// compress bytes to buffer
func GzipToBuffer(data []byte, buf *bytes.Buffer) error {
	wr := gzip.NewWriter(buf)
	_, err := wr.Write(data)
	if err != nil {
		return err
	}
	err = wr.Close()
	if err != nil {
		return err
	}
	return nil
}
