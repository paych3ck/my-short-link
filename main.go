package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DbUser     string `json:"dbUser"`
	DbPassword string `json:"dbPassword"`
	DbHost     string `json:"dbHost"`
	DbName     string `json:"dbName"`
}

func generatePasswordHash(password string) string {
	passwordBytes := []byte(password)
	cost := 10
	hash, _ := bcrypt.GenerateFromPassword(passwordBytes, cost)
	return string(hash)
}

func checkPasswordHash(password string, hashedPassword string) bool {
	passwordBytes := []byte(password)
	hashedPasswordBytes := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes)
	return err == nil
}

func connectDB() *sql.DB {
	file, err := os.Open("config.json")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)

	if err != nil {
		log.Fatal(err)
	}

	dsn := config.DbUser + ":" + config.DbPassword + "@tcp(" + config.DbHost + ")/" + config.DbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func createBind(db *sql.DB, alias string, url string) error {
	_, err := db.Exec("INSERT INTO urls(alias, url) VALUES(?, ?)", alias, url)
	if err != nil {
		return err
	}

	return nil
}

func registerUser(db *sql.DB, email, password string) error {
	_, err := db.Exec("INSERT INTO users(email, password) VALUES(?, ?)", email, password)

	if err != nil {
		return err
	}

	return nil
}

func sendEmail(to, subject, body string) {
	//TODO: SEND EMAIL HERE
}

func checkUser(db *sql.DB, email, password string) (bool, error) {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&dbPassword)

	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	return checkPasswordHash(password, dbPassword), nil
}

func generateShortUrl(length int) string {
	const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)

	for i := range bytes {
		bytes[i] = symbols[rand.Intn(len(symbols))]
	}

	return "https://myshl.ru/" + string(bytes)
}

func main() {
	fs := http.FileServer(http.Dir("static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("index.html")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		db := connectDB()
		defer db.Close()
		shortUrl := generateShortUrl(3)
		createBind(db, shortUrl, requestData.URL)
		json.NewEncoder(w).Encode(map[string]string{"shortUrl": shortUrl})
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("login.html")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {
			email := r.FormValue("email")
			password := r.FormValue("password")

			db := connectDB()
			defer db.Close()

			isValid, err := checkUser(db, email, password)

			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if isValid {
				fmt.Println("succes")
			} else {
				fmt.Println("not succes")
			}
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			tmpl, err := template.ParseFiles("register.html")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, nil)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else if r.Method == "POST" {
			email := r.FormValue("email")
			password := r.FormValue("password")

			db := connectDB()
			defer db.Close()

			registerUser(db, email, generatePasswordHash(password))
		}
	})

	http.ListenAndServe(":8080", nil)
}
