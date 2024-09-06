package main

import (
	"encoding/json"
	"fmt"
	"internal/adapters/logger"
	"internal/app"
	"internal/domain"
	"net/http"
	"time"
)

func submit(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	if err := r.ParseForm(); err != nil {
		logger.Error(fmt.Sprintf("Quiz parse error: %s\n", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	qrp := app.GetReplies(r.Form)

	jQrp, err := json.Marshal(qrp)
	if err != nil {
		logger.Error(fmt.Sprintf("Quiz parse error: %s\n", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := r.Form.Get("test_id")
	if uuid == "" {
		logger.Error(fmt.Sprintf("Submit: Wrong test ID [%s]\n", uuid))
		http.Error(w, "Error. Unexpected quiz ID", http.StatusBadRequest)
		return
	}

	qd, err := stor.GetQuizData(uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("Quiz extract error [%s]: %s\n", uuid, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Quiz not found [%s]", uuid), http.StatusInternalServerError)
		return
	}

	qst, err := stor.GetQuizScores(uuid)
	if err != nil {
		logger.Error(fmt.Sprintf("Quiz scores extract error [%s]: %s\n", uuid, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Quiz not found [%s]", uuid), http.StatusInternalServerError)
		return
	}
	qd.Scores = append(qd.Scores, *qst...)

	//TODO: move more of quiz logic to app
	totalCorrrect := app.GetQuizTotalScore(qd)
	totalAnswered := app.GetQuizUserScore(qd, qrp)

	if totalCorrrect == 0 {
		logger.Error(fmt.Sprintf("Zero correct answers in quiz [%s]\n", uuid))
		http.Error(w, fmt.Sprintf("Error. Quiz result calculation error [%s]", uuid), http.StatusInternalServerError)
		return
	}

	totalPercent := app.GetQuizUserRatio(totalAnswered, totalCorrrect)

	qs, err := app.GetQuizScore(qd, totalPercent)
	if err != nil {
		logger.Error(fmt.Sprintf("Quiz calculation error [%s]: %s\n", uuid, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Quiz result calculation error [%s]", uuid), http.StatusInternalServerError)
		return
	}

	var qr domain.QuizResult
	qr.TestID = uuid
	qr.ScoreID = qs.UUID
	qr.PassTime = time.Now()
	qr.Result = totalPercent
	qr.Score = qs.Score
	qr.IsPassed = qs.Pass
	qr.Replies = string(jQrp)

	logger.Info(fmt.Sprintf("RESULTS: %2d/%2d (ratio %d%%)\n", totalAnswered, totalCorrrect, totalPercent)) //DEBUG
	logger.Info(fmt.Sprintf("RESULTS: %v\n", qr.Replies))                                                   //DEBUG

	resultID, err := stor.WriteQuizResult(qr)
	if err != nil || resultID == "" {
		logger.Error(fmt.Sprintf("Quiz result write error [%s]: %s\n", uuid, err.Error()))
		http.Error(w, fmt.Sprintf("Error. Quiz result write error [%s]", uuid), http.StatusInternalServerError)
		return
	}

	urlID := domain.EncodeGUID(resultID)

	http.Redirect(w, r, fmt.Sprintf("/results/%s", urlID), http.StatusSeeOther)

}
