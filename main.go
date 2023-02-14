package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"example.com/webservice/data"
	_ "modernc.org/sqlite"
)

// GET  /user/1  -> Retrieve User Id 1
// PUT  /user/1  -> Update User Id 1
// POST /user    -> Create new user, return user record
// DELETE /user/1 -> Delete User Id 1

func main() {
	http.HandleFunc("/hello", SayHello)
	http.HandleFunc("/user/", UserDispatch)

	http.ListenAndServe(":8080", nil)
}

func UserDispatch(rw http.ResponseWriter, r *http.Request) {

	PathStart := "/user/"
	_, pathinfo, _ := strings.Cut(r.URL.Path, PathStart)

	userId := 0
	userId, err := strconv.Atoi(pathinfo)
	if err != nil {
		userId = 0
	}

	switch r.Method {
	case "GET":
		if pathinfo == "" {
			ListUsers(rw, r)
			return
		} else {
			ShowUser(rw, r, userId)
		}
		break
	case "POST":
		if pathinfo == "" {
			CreateUser(rw, r)
			return
		}
		break
	default:
		ListUsers(rw, r)
		break
	}

}

func CreateUser(rw http.ResponseWriter, r *http.Request) {

	var u data.User

	body, err := io.ReadAll(r.Body)
	if len(body) == 0 {
		rw.WriteHeader(400)
		return
	}

	json.Unmarshal(body, &u)

	db, err := data.Connect()
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	result, err := db.Exec("INSERT INTO user (first_name, last_name, email) VALUES ($1, $2, $3)", u.FirstName, u.LastName, u.Email)
	id, err := result.LastInsertId()
	u.Id = id

	output, err := json.MarshalIndent(u, "", "  ")
	rw.Write(output)

}

func ShowUser(rw http.ResponseWriter, r *http.Request, userid int) {
	db, err := data.Connect()
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rows := db.QueryRow("SELECT * FROM user WHERE id=$1", userid)

	var u data.User
	rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email)

	output, err := json.Marshal(u)
	rw.Write(output)

}

func ListUsers(rw http.ResponseWriter, r *http.Request) {
	db, err := data.Connect()
	if err != nil {
		rw.WriteHeader(500)
		return
	}

	rows, err := db.Query("SELECT * FROM user;")
	if err != nil {
		log.Printf("Retrieving users: %v", err)
		rw.WriteHeader(500)
		return
	}
	defer rows.Close()

	var users []data.User
	for rows.Next() {
		var u data.User

		if err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email); err != nil {
			log.Fatal(err)
		}

		users = append(users, u)
	}

	output, err := json.MarshalIndent(users, "", "  ")
	rw.Write(output)
}

func SayHello(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Parsing form: %v", err)
		rw.WriteHeader(400)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		name = "Friend"
	}
	rw.Write([]byte(fmt.Sprintf("Hello %s", name)))
}
