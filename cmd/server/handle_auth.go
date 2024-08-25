package main

import (
	"fmt"
	"internal/adapters/logger"
	"net/http"
	"time"
)

func auth(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	ip := getClientIP(r)
	key := ip + ":" + username

	logger.Info(fmt.Sprintf("Auth attempt: %v\n", key))

	if isUSerBlocked(key) {
		http.Error(w, "Вы заблокированы из-за слишком большого количества неудачных попыток входа. Попробуйте позже.", http.StatusTooManyRequests)
		return
	}

	if username == adminUser && checkPasswordHash(password, adminPassword) {
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "authenticated",
			Path:     "/",
			HttpOnly: true,
			//Secure:   true, Not possible while debugging
			Expires: time.Now().Add(24 * time.Hour), // TTL 1 день
		})
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	registerFailedAuth(key)

	http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
}

func auth_page(w http.ResponseWriter, r *http.Request) {

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`
	<html>
	<style>
		body {
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
			font-family: Arial, sans-serif;
		}
		h2 {
			text-align: center;
		}
	</style>
	<body>
		<div class="login-container">
			<h2>Введите данные для авторизации</h2>
			<form method="POST">
				<label for="username">Логин:</label><br>
				<input type="text" id="username" name="username" required><br>
				<label for="password">Пароль:</label><br>
				<input type="password" id="password" name="password" required><br>
				<button type="submit">Войти</button>
			</form>
		</div>
	</body>
	</html>
`))

}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
