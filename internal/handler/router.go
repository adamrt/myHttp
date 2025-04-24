package handler

import (
	"errors"
	"fmt"
)

type Router struct {
	Routes map[string]map[string]handler
}

func (r *Router) Handle(method, target string) (handler, error) {
	fmt.Printf("Router method called, Method:%s, Target:%s", method, target)
	if elem, ok := r.Routes[method][target]; ok {
		fmt.Printf("Handler function recieved: %v", elem)
		return elem, nil
	}
	return nil, errors.New("Handler doesn't exist")
}
func (r *Router) SetRouters(method, target string, handle handler) {
	if _, ok := r.Routes[method]; !ok {
		r.Routes[method] = make(map[string]handler)
	}
	r.Routes[method][target] = handle
}

func NewRouter() *Router {
	return &Router{
		Routes: make(map[string]map[string]handler),
	}
}
