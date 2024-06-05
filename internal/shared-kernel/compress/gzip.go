package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
)

type Writer struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *Writer {
	return &Writer{
		w,
		gzip.NewWriter(w),
	}
}

func (c *Writer) Write(p []byte) (int, error) {
	n, err := c.zw.Write(p)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}
	return n, nil
}

func (c *Writer) Close() error {
	err := c.zw.Close()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func GzipData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return []byte{}, fmt.Errorf("%w", err)
	}
	err = gz.Close()
	if err != nil {
		return []byte{}, fmt.Errorf("%w", err)
	}
	return b.Bytes(), nil
}
