package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int
	Name     string
	DOB      string
	Gender   string
	Username string
	Password string
	Balance  float64
}

type Transaction struct {
	Type      string
	Amount    float64
	Timestamp time.Time
}

var (
	db   *sql.DB
	tmpl *template.Template
)

const (
	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultDBUser     = "postgres"
	defaultDBPassword = "Postgresql"
	defaultDBName     = "bankapp"
	defaultAppPort    = "8080"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func initDB() {
	host := getEnv("DB_HOST", defaultDBHost)
	port := getEnv("DB_PORT", defaultDBPort)
	user := getEnv("DB_USER", defaultDBUser)
	password := getEnv("DB_PASSWORD", defaultDBPassword)
	dbname := getEnv("DB_NAME", defaultDBName)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
}

func initTemplates() {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl.ExecuteTemplate(w, "register.html", nil)
	case http.MethodPost:
		name := r.FormValue("name")
		dob := r.FormValue("dob")
		gender := r.FormValue("gender")
		username := r.FormValue("username")
		password := r.FormValue("password")
		confirm := r.FormValue("confirm")

		if password != confirm {
			tmpl.ExecuteTemplate(w, "register.html", "Passwords do not match")
			return
		}

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			tmpl.ExecuteTemplate(w, "register.html", "Username already exists")
			return
		}

		_, err = db.Exec("INSERT INTO users (name, dob, gender, username, password, balance) VALUES ($1, $2, $3, $4, $5, 0.0)",
			name, dob, gender, username, password)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl.ExecuteTemplate(w, "login.html", nil)
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbUser User
		err := db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).
			Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				tmpl.ExecuteTemplate(w, "login.html", "Invalid credentials")
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		if password == dbUser.Password {
			http.Redirect(w, r, "/dashboard?user="+username, http.StatusSeeOther)
		} else {
			tmpl.ExecuteTemplate(w, "login.html", "Invalid credentials")
		}
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")

	var user User
	err := db.QueryRow("SELECT id, name, dob, gender, username, balance FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Name, &user.DOB, &user.Gender, &user.Username, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	data := struct {
		User  *User
		Error string
	}{
		User:  &user,
		Error: r.URL.Query().Get("error"),
	}

	tmpl.ExecuteTemplate(w, "dashboard.html", data)
}

func depositHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	amountStr := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Redirect(w, r, "/dashboard?user="+username+"&error=Invalid+amount", http.StatusSeeOther)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE username = $2", amount, username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var userID int
	err = tx.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO transactions (user_id, type, amount) VALUES ($1, $2, $3)",
		userID, "Deposit", amount)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard?user="+username, http.StatusSeeOther)
}

func withdrawHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	amountStr := r.FormValue("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Redirect(w, r, "/dashboard?user="+username+"&error=Invalid+amount", http.StatusSeeOther)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var balance float64
	err = tx.QueryRow("SELECT balance FROM users WHERE username = $1", username).Scan(&balance)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if balance < amount {
		http.Redirect(w, r, "/dashboard?user="+username+"&error=Insufficient+funds", http.StatusSeeOther)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE username = $2", amount, username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var userID int
	err = tx.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO transactions (user_id, type, amount) VALUES ($1, $2, $3)",
		userID, "Withdrawal", amount)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard?user="+username, http.StatusSeeOther)
}

func transactionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")

	var user User
	err := db.QueryRow("SELECT id, name, username, balance FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Name, &user.Username, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	rows, err := db.Query("SELECT type, amount, timestamp FROM transactions WHERE user_id = $1 ORDER BY timestamp DESC", user.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.Type, &t.Amount, &t.Timestamp)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}

	data := struct {
		User         *User
		Transactions []Transaction
	}{
		User:         &user,
		Transactions: transactions,
	}

	tmpl.ExecuteTemplate(w, "history.html", data)
}

func main() {
	initDB()
	initTemplates()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/deposit", depositHandler)
	http.HandleFunc("/withdraw", withdrawHandler)
	http.HandleFunc("/history", transactionHistoryHandler)

	appPort := getEnv("APP_PORT", defaultAppPort)
	fmt.Printf("Server running at http://localhost:%s/\n", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, nil))
}
