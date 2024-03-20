package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DbUser     string `json:"dbUser"`
	DbPassword string `json:"dbPassword"`
	DbHost     string `json:"dbHost"`
	DbName     string `json:"dbName"`
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

func registerUser(db *sql.DB, email, password string) error {
	_, err := db.Exec("INSERT INTO users(email, password) VALUES(?, ?)", email, password)

	if err != nil {
		return err
	}

	return nil
}

func checkUser(db *sql.DB, email, password string) (bool, error) {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&dbPassword)

	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	return password == dbPassword, nil
}

func main() {
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
		} else if r.Method == "POST" {
			fmt.Println("test")
		}
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

			registerUser(db, email, password)
		}
	})

	http.ListenAndServe(":8080", nil)
}
