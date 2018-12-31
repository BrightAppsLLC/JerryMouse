package clients

import (
	"honnef.co/go/js/dom"
)

// ClientAppContext -
type ClientAppContext struct {
	TemplateSet          *isokit.TemplateSet
	Window               dom.Window
	Document             dom.Document
	PageContentContainer dom.Element
}
