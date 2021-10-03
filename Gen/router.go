package gen

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*trieNode
	handlers map[string]HandlerFunc
}

//roots key eg, roots['GET'] roots['POST']
//handlers key eg, handlers['GET-/p/:lang/doc], handlers['POST-/p/book']

func newRouter() *router {
	return &router{
		roots:    make(map[string]*trieNode),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			//the other parts after "*" won't be save
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	if _, has := r.roots[method]; !has {
		r.roots[method] = &trieNode{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*trieNode, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, has := r.roots[method]
	if !has {
		return nil, nil
	}
	node := root.search(searchParts, 0)

	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] //??
			}
			if part[0] == '*' {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)

	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s/n", c.Path)
		})
	}

	//do the middlewares and handlers
	c.Next()
}
