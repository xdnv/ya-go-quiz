package main

import (
	"fmt"
	"internal/adapters/logger"
	"io"
	"net/http"
)

func handleCommand(w http.ResponseWriter, r *http.Request) {

	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	// Печатаем данные в консоль
	logger.Info(fmt.Sprintf("Got data: %s\n", body))

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusOK)
	//return
}
