package utils

import (
	"html/template"
)

// Sink is a struct that identifies the name
// of the sink, the associated URL and the
// HTTP method
type Sink struct {
	Name   template.HTML
	URL    template.HTML
	Method template.HTML
}

// Route is the template information for a specific route
type Route struct {
	Base     template.HTML
	Name     template.HTML
	Link     template.HTML
	Products []template.HTML
	Inputs   []template.HTML
	Sinks    []Sink
}

// Parameters are the parameters for a specific page
type Parameters struct {
	Body    template.HTML
	Year    int
	Rulebar map[string]Route
	Port    string
	Name    string
}
