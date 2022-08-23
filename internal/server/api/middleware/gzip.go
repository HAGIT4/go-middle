// Package middleware Middlewares for collector service
package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

func GzipWriteMiddleware() (h gin.HandlerFunc) {
	return func(c *gin.Context) {
		clientAcceptEncodingHeaders := c.Request.Header.Values("Accept-Encoding")
		if len(clientAcceptEncodingHeaders) == 0 {
			return
		}
		clientAcceptEncodingHeadersSlice := strings.Split(clientAcceptEncodingHeaders[0], ",")
		clientAcceptEncodingHeadersMap := make(map[string]bool)
		for _, header := range clientAcceptEncodingHeadersSlice {
			clientAcceptEncodingHeadersMap[header] = true
		}

		if _, found := clientAcceptEncodingHeadersMap["gzip"]; !found {
			return
		}

		c.Header("Content-Encoding", "gzip")
		gz := gzip.NewWriter(c.Writer)
		c.Writer = &gzipWriter{
			c.Writer,
			gz,
		}

		defer func() {
			gz.Close()
			c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
		}()
		c.Next()
	}
}

func GzipReadMiddleware() (h gin.HandlerFunc) {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(c.Request.Body)
			if err == io.EOF {
				return
			} else if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			defer gz.Close()

			body, err := io.ReadAll(gz)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		}
	}
}
