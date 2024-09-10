package main

import (
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func handleCommand(w http.ResponseWriter, r *http.Request) {

	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusOK)

	command := chi.URLParam(r, "command")
	url := chi.URLParam(r, "id")

	logger.Info(fmt.Sprintf("Command: %s, id: %s\n", command, url)) //DEBUG

	uuid, err := domain.DecodeGUID(url)
	if err != nil {
		logger.Error(fmt.Sprintf("Wrong Quiz ID [%s]: %s\n", url, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	if strings.TrimSpace(uuid) == "" {
		logger.Error(fmt.Sprintf("Wrong Quiz ID [%s]\n", url))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Error reading body: %v\n", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	logger.Info(fmt.Sprintf("Body: %s\n", body))

	switch command {
	case "toggle":
		err = stor.ToggleQuizAvailability(uuid)
		if err != nil {
			logger.Error(fmt.Sprintf("Wrong Quiz ID [%s]: %s\n", url, err.Error()))
			http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
			return
		}
	default:
		logger.Error(fmt.Sprintf("Unknown command: %s\n", command))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", command), http.StatusNotFound)
		return
	}
}
