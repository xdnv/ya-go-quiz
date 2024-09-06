package main

import (
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handleCommand(w http.ResponseWriter, r *http.Request) {

	logger.Info("hello")

	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	logger.Info("hello2")

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusOK)

	command := chi.URLParam(r, "command")
	url := chi.URLParam(r, "id")
	//logger.Info(fmt.Sprintf("Command: %s\n", command))

	if command != "toggle" {
		logger.Error(fmt.Sprintf("Wrong command: %s\n", command))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", command), http.StatusNotFound)
		return
	}

	uuid, err := domain.DecodeGUID(url)
	if err != nil {
		logger.Error(fmt.Sprintf("Wrong Quiz ID [%s]: %s\n", url, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

	if uuid == "" {
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

	logger.Info(fmt.Sprintf("Got data: %s\n", body))

	err = stor.ToggleQuizAvailability(uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("Wrong Quiz ID [%s]: %s\n", url, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Resource not found: [%s]", url), http.StatusNotFound)
		return
	}

}
