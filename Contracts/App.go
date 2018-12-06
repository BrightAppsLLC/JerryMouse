package Contracts

import (
	"fmt"
	"net/http"

	"github.com/TheFlyingNomad/JerryMouse/Contracts"
)

// RenderTemplate -
func RenderTemplate(w http.ResponseWriter, appContext *Contracts.ServerAppContext, templateFile string, templateData interface{}) {
	executeError := appContext.Templates.ExecuteTemplate(w, templateFile, templateData)
	if executeError != nil {
		fmt.Printf("RenderTemplate: %s", executeError)
	}
}
