package gen

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

//print the stack message for debug
func trace(msg string) string {
	var pcs [32]uintptr
	//0 caller is Callers itself, 1 caller is trace, 2 caller is defer func()
	n := runtime.Callers(3, pcs[:]) //Skip First 3 Callers

	var str strings.Builder
	str.WriteString(msg + "\nTraceBack:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

//recovery function, prevent sever from crash down
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(msg))
				c.Fail(http.StatusInternalServerError, errors.New("Internal Server Error"))
			}
		}()

		c.Next()
	}
}
