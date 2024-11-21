package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func writeHTMLResponse(w http.ResponseWriter, html template.HTML) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleError(w http.ResponseWriter, operation string, err error, code int) {
	log.Printf("Error during %s: %v", operation, err)
	http.Error(w, fmt.Sprintf("Error during %s", operation), code)
}
