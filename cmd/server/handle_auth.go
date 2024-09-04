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
		http.Error(w, "Превышено количество неудачных попыток входа. Попробуйте позже.", http.StatusTooManyRequests)
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

func authPage(w http.ResponseWriter, r *http.Request) {

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
		.login-container {
			border: 1px solid #ccc;
			padding: 20px;
			border-radius: 5px;
			box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
			width: 300px;
		}
		h2 {
			text-align: center;
		}
		label {
			display: block;
			margin-bottom: 5px;
		}
		input[type="text"],
		input[type="password"] {
			width: 100%;
			padding: 10px;
			margin-bottom: 15px;
			border: 1px solid #ccc;
			border-radius: 4px;
		}
		button {
			width: 100%;
			padding: 10px;
			background-color: #28a745;
			color: white;
			border: none;
			border-radius: 4px;
			cursor: pointer;
		}
		button:hover {
			background-color: #218838;
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
