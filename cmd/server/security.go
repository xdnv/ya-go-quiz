package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type BlockedUser struct {
	FailedAttempts int
	LastAttemptAt  time.Time
}

var (
	//simplified auth engine. Yes, it should be stored in database
	adminUser     = "admin"
	adminPassword = "$2a$10$/w.tr1Hd1TCzBHkDC2nmCOfThwIQseOG/K/EEowuv44ii5XNDefUe" // "password" hash
	blockedUsers  = make(map[string]BlockedUser)
	mu            sync.Mutex
)

// Hash password provided by user
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare password with hashed one
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	//logger.Error(fmt.Sprintf("pwd compare error: %v\n", err)) //DEBUG
	return err == nil
}

// Get client IP-address
func getClientIP(r *http.Request) string {
	//TODO: strip port from return value
	// check X-Forwarded-For header if we sit behing proxy/balancer
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// Get first IP if there's more than one
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// Use RemoteAddr if there's no header
	ip = r.RemoteAddr
	return ip
}

func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	return err == nil && cookie.Value == "authenticated"
}

func isUSerBlocked(username string) bool {
	result := false
	mu.Lock()
	if blockedUser, ok := blockedUsers[username]; ok {
		//Block period is still active
		blockDuration := time.Minute
		if blockedUser.FailedAttempts > 3 {
			blockDuration = time.Hour
		}

		if time.Since(blockedUser.LastAttemptAt) < time.Duration(blockedUser.FailedAttempts*blockedUser.FailedAttempts)*blockDuration {
			result = true
		} else {
			//Block period has ended, release user+ip from prison
			delete(blockedUsers, username)
		}
	}
	mu.Unlock()
	return result
}

func registerFailedAuth(username string) {
	mu.Lock()
	if blockedUser, ok := blockedUsers[username]; ok {
		blockedUser.FailedAttempts++
		blockedUser.LastAttemptAt = time.Now()
	} else {
		blockedUsers[username] = BlockedUser{
			FailedAttempts: 1,
			LastAttemptAt:  time.Now(),
		}
	}
	mu.Unlock()
}
