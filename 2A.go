//go:build !solution

package otp

import (
	"io"
)

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &streamReader{r, prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &streamWriter{w, prng}
}

type streamReader struct {
	r    io.Reader
	prng io.Reader
}

func (sr *streamReader) Read(p []byte) (n int, err error) {
	n, err = sr.r.Read(p)

	if err != nil {
		// If the error is io.EOF, we still need to XOR the bytes that were read
		if err == io.EOF {
			prngBuf := make([]byte, n)
			sr.prng.Read(prngBuf)
			for i := 0; i < n; i++ {
				p[i] ^= prngBuf[i]
			}
			return n, io.EOF
		}

		return n, err
	}

	prngBuf := make([]byte, n)
	sr.prng.Read(prngBuf)

	for i := 0; i < n; i++ {
		p[i] ^= prngBuf[i]
	}

	return
}

type streamWriter struct {
	w    io.Writer
	prng io.Reader
}

func (sw *streamWriter) Write(p []byte) (n int, err error) {
	prngBuf := make([]byte, len(p))
	sw.prng.Read(prngBuf)

	buf := make([]byte, len(p))
	copy(buf, p)

	for i := 0; i < len(p); i++ {
		buf[i] ^= prngBuf[i]
	}

	return sw.w.Write(buf)
}
