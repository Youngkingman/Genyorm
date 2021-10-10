package gen

import (
	"errors"
	"net/http"
)

// 处理跨域请求,支持options访问
func Cors() HandlerFunc {
	return func(c *Context) {
		method := c.Req.Method

		c.SetHeader("Access-Control-Allow-Origin", "*")
		c.SetHeader("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.SetHeader("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.SetHeader("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.JSON(http.StatusNoContent, errors.New("NoContent"))
		}

		c.Next()
	}
}
