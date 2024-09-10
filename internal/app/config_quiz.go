package app

import (
	"fmt"
	"internal/domain"
	"net/url"
	"strings"
)

func GetReplies(rf url.Values) *domain.QuizReplies {
	qr := make(domain.QuizReplies)
	for key, values := range rf {
		qr[key] = values
	}
	return &qr
}

// calculate max/total score
func GetQuizTotalScore(qd *domain.QuizData) int {

	totalCorrrect := 0

	for _, q := range qd.Questions {
		switch q.Type {
		case "single_choice":
			for _, o := range q.Options {
				if o.IsCorrect {
					totalCorrrect++
				}
			}
		case "multiple_choice":
			for _, o := range q.Options {
				if o.IsCorrect {
					totalCorrrect++
				}
			}
		case "entry_int":
			totalCorrrect++
		}
	}
	return totalCorrrect
}

// calculate max/total score
func GetQuizUserScore(qd *domain.QuizData, rf *domain.QuizReplies) int {

	totalAnswered := 0

	//calculate actual score
	for _, q := range qd.Questions {
		if len((*rf)[q.UUID]) == 0 {
			continue
		}
		result := (*rf)[q.UUID][0]
		if strings.TrimSpace(result) == "" {
			continue
		}
		switch q.Type {
		case "single_choice":
			for _, o := range q.Options {
				if result == o.UUID && o.IsCorrect {
					//logger.Info(fmt.Sprintf("+1 %s/%s\n", q.Text, o.Text)) //DEBUG
					totalAnswered++
				}
			}
		case "multiple_choice":
			for _, o := range q.Options {
				for _, value := range (*rf)[q.UUID] {
					if value == o.UUID && o.IsCorrect {
						//logger.Info(fmt.Sprintf("+1 %s/%s\n", q.Text, o.Text)) //DEBUG
						totalAnswered++
					}
				}
			}
		case "entry_int":
			if result == q.Options[0].Value {
				//logger.Info(fmt.Sprintf("+1 %s/%s\n", q.Text, q.Options[0].Text)) //DEBUG
				totalAnswered++
			}
		}
	}
	return totalAnswered
}

// calculate success ratio in decimal %
func GetQuizUserRatio(answered int, corrrect int) int {
	//int precent calculation results in either 0 or 100%
	fa := float64(answered)
	fc := float64(corrrect)
	return int((fa / fc) * 100)
}

// returns existing ScoreData object corresponding to the ratio provided
func GetQuizScore(qd *domain.QuizData, ratio int) (*domain.QuizScore, error) {
	for _, s := range qd.Scores {
		if ratio >= s.MinPercent && ratio <= s.MaxPercent {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("no quizscore objects reflect ratio %d", ratio)
}
