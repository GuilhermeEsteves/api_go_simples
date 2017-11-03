package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func main() {
	configureConnectionDB()
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/getUsers", getUsers)
	http.HandleFunc("/postUser", postUser)
	fmt.Println("Api rodando na porta 1200")
	log.Fatal(http.ListenAndServe(":1200", nil))
}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Nao autorizado", 401)
		return
	}
	// users := []User{
	// 	User{
	// 		Name: "Esteves",
	// 		Age:  22,
	// 	},
	// 	User{
	// 		Name: "Gerepe",
	// 		Age:  22,
	// 	},
	// }

	users := []User{}

	rows, err := db.Query("Select * from User")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		user := User{}

		err := rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			fmt.Println(err)
		}

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func postUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Nao autorizado", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	user := User{}

	errDecode := decoder.Decode(&user)
	if errDecode != nil {
		fmt.Println("Error decoder", errDecode)
		return
	}

	query := "Insert into User values(?,?,?)"

	_, err := db.Exec(query, user.Id, user.Name, user.Age)
	if err != nil {
		fmt.Println("Error", err)
	}

	w.WriteHeader(http.StatusCreated)
}

type User struct {
	Id   int
	Name string
	Age  int
}

func configureConnectionDB() {
	var errorBanco error
	db, errorBanco = sql.Open("mysql", "root:123@tcp(localhost:3306)/go")
	if errorBanco != nil {
		fmt.Println(errorBanco)
	}

	err := db.Ping()
	if err != nil {
		fmt.Println("Erro no ping do banco")
	}

	fmt.Println("Connect Mysql :)")
}
