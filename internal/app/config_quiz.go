package app

import (
	"internal/domain"
	"net/url"
)

// calculate max/total score
func GetQuizTotalScore(qd *domain.QuizData) int {

	total_corrrect := 0

	for _, q := range qd.Questions {
		switch q.Type {
		case "single_choice":
			for _, o := range q.Options {
				if o.IsCorrect {
					total_corrrect++
				}
			}
		case "multiple_choice":
			for _, o := range q.Options {
				if o.IsCorrect {
					total_corrrect++
				}
			}
		case "entry_int":
			total_corrrect++
		}
	}
	return total_corrrect
}

// calculate max/total score
func GetQuizUserScore(qd *domain.QuizData, rf url.Values) int {

	total_answered := 0

	//calculate actual score
	for _, q := range qd.Questions {
		result := rf.Get(q.UUID)
		if result == "" {
			continue
		}
		switch q.Type {
		case "single_choice":
			for _, o := range q.Options {
				if result == o.UUID && o.IsCorrect {
					//logger.Info(fmt.Sprintf("+1 %s/%s\n", q.Text, o.Text)) //DEBUG
					total_answered++
				}
			}
		case "multiple_choice":
			for _, o := range q.Options {
				for _, value := range rf[q.UUID] {
					if value == o.UUID && o.IsCorrect {
						//logger.Info(fmt.Sprintf("+1 %s/%s\n", q.Text, o.Text)) //DEBUG
						total_answered++
					}
				}
			}
		case "entry_int":
			if result == q.Options[0].Value {
				//logger.Info(fmt.Sprintf("+1 %s/%s\n", q.Text, q.Options[0].Text)) //DEBUG
				total_answered++
			}
		}
	}
	return total_answered
}

// calculate success ratio in decimal %
func GetQuizUserRatio(total_answered int, total_corrrect int) int {
	//int precent calculation results in either 0 or 100%
	fta := float64(total_answered)
	ftc := float64(total_corrrect)
	return int((fta / ftc) * 100)
}
