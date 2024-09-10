package main

import (
	"encoding/json"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"net/http"
	"strings"
)

func uploadData(w http.ResponseWriter, r *http.Request) {

	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var maxMemory int64 = sc.MaxFileMemory
	var qd domain.QuizData
	var dec *json.Decoder

	//enforce maxMemory MB body size, else --> err="http: request body too large"
	r.Body = http.MaxBytesReader(w, r.Body, maxMemory*1024*1024)

	//Multipart Header:
	//Content-Type: multipart/form-data; boundary=----WebKitFormBoundary44A4qJbyQuzKMjvB

	//Multipart Body:
	//Content-Disposition: form-data; name="file"; filename="demoTest001.json"
	//Content-Type: application/json

	ct := r.Header.Get("Content-Type")
	if ct == "" {
		msg := "Content-Type header must be set"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	if mediaType == "multipart/form-data" {
		err := r.ParseMultipartForm(maxMemory << 20)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to parse multipart message: %s\n", err.Error()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get file from multipart object
		file, header, err := r.FormFile("file")
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get multipart file: %s\n", err.Error()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Print multipart file headers
		logger.Info(fmt.Sprintf("Uploaded File: %s\n", header.Filename))
		logger.Info("Headers:\n")
		for key, value := range header.Header {
			logger.Info(fmt.Sprintf("%s: %s\n", key, value))
		}
		dec = json.NewDecoder(file)
	} else if mediaType == "application/json" {
		dec = json.NewDecoder(r.Body)
	} else {
		msg := "Content-Type header should be application/json or multipart/form-data"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	//force err="json: unknown field ..."
	dec.DisallowUnknownFields()

	if err := dec.Decode(&qd); err != nil {
		logger.Error(fmt.Sprintf("Error decoding Quiz config: %s\n", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info(fmt.Sprintf("Got quiz: %s v. %s\n", qd.ID, qd.Version)) //DEBUG

	var errs []error

	stor.UpdateQuiz(&qd, &errs)

	//handling all errors encountered
	if len(errs) > 0 {
		strErrors := make([]string, len(errs))
		for i, err := range errs {
			strErrors[i] = err.Error()
		}
		errDesc := strings.Join(strErrors, "\n")

		logger.Error(fmt.Sprintf("Error saving Quiz config: %s\n", errDesc))
		http.Error(w, "Errors: \n"+errDesc, http.StatusBadRequest)
		return
	}

	//w.WriteHeader(http.StatusOK)
}
