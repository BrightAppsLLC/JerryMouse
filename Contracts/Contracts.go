package Contracts

import (
	"html/template"

	"github.com/IsomorphicGo/isokit"
	"honnef.co/go/js/dom"
)

// EmptyObject -
type EmptyObject struct{}

// ClientAppContext -
type ClientAppContext struct {
	TemplateSet          *isokit.TemplateSet
	Window               dom.Window
	Document             dom.Document
	PageContentContainer dom.Element
}

// ServerAppContext -
type ServerAppContext struct {
	Templates   *template.Template
	TemplateSet *isokit.TemplateSet
}
