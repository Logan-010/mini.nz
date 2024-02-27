package lib

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		handler.ServeHTTP(gzw, r)
	})
}

func compress(data []byte) ([]byte, error) {
	var compressedData bytes.Buffer
	gzw := gzip.NewWriter(&compressedData)

	_, err := gzw.Write(data)
	if err != nil {
		return nil, err
	}

	err = gzw.Close()
	if err != nil {
		return nil, err
	}

	return compressedData.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(data)

	gzr, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	decompressed, err := io.ReadAll(gzr)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}
