package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rd1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rd2 := ioutil.NopCloser(bytes.NewBuffer(buf))

		fmt.Println("Req Body:\n", readBody(rd1))

		c.Request.Body = rd2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
