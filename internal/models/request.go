package models

import "fmt"

type Request struct {
	Name   string
	Method string
	Target Target
	Path   string
}

func (r *Request) Run() {
	fmt.Printf("%s %s %s %s\n", r.Name, r.Method, r.Target.URL, r.Path)
}
