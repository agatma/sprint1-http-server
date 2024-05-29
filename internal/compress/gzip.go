package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

type Writer struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *Writer {
	return &Writer{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *Writer) Header() http.Header {
	return c.w.Header()
}

func (c *Writer) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *Writer) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *Writer) Close() error {
	err := c.zw.Close()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

type Reader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*Reader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Reader{
		r:  r,
		zr: zr,
	}, nil
}

func (c Reader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *Reader) Close() error {
	if err := c.r.Close(); err != nil {
		return fmt.Errorf("%w", err)
	}
	return c.zr.Close()
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
