package main

import (
	gen "Ghenyorm/Gen"
	"net/http"
)

func main() {
	r := gen.New()
	r.GET("/hello", func(c *gen.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	//pack the form value into json and write back to requester
	r.POST("/login", func(c *gen.Context) {
		c.JSON(http.StatusOK, gen.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run(":1234")
}
