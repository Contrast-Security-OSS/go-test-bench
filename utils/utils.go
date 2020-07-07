package utils

import (
	"html/template"
)

type Sink struct {
	Name   template.HTML
	Url    template.HTML
	Method template.HTML
}

type Route struct {
	Base     template.HTML
	Name     template.HTML
	Link     template.HTML
	Products []template.HTML
	Inputs   []template.HTML
	Sinks    []Sink
}

type Parameters struct {
	Body    template.HTML
	Year    int
	Rulebar map[string]Route
	Port    string
	Name    string
}
