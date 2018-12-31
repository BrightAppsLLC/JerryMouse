package clients

import (
	"fmt"
	"net/http"
)

// RenderTemplate -
func RenderTemplate(w http.ResponseWriter, appContext *ServerAppContext, templateFile string, templateData interface{}) {
	executeError := appContext.Templates.ExecuteTemplate(w, templateFile, templateData)
	if executeError != nil {
		fmt.Printf("RenderTemplate: %s", executeError)
	}
}
