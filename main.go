package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Psw   string `json:"psw"`
}
type Post struct {
	ID         uint   `json:"id"`
	Text       string `json:"text"`
	User_id    uint   `json:"user_id"`
	User_login string `json:"login"`
}
type LoginStatus struct {
	Login  string `json:"login"`
	Status int    `json:"status"`
}

var db *sql.DB
var err error

func GetUserPosts(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	req_params := mux.Vars(req)
	var rows *sql.Rows
	u_id := req_params["login"]
	fmt.Println("Get posts")
	fmt.Println(u_id)
	fmt.Println("select post.id, post.text, post.user_id, blog_user.login from post join blog_user where login = '" + u_id + "'")
	if u_id == "default" {
		fmt.Println("all")
		rows, err = db.Query("select post.id, post.text, post.user_id, blog_user.login from post join blog_user on post.user_id = blog_user.id")
	} else {
		fmt.Println("not all")
		rows, err = db.Query("select post.id, post.text, post.user_id, blog_user.login from post join blog_user on post.user_id = blog_user.id where login = '" + u_id + "'")
	}
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var posts []Post
	p := Post{}

	for rows.Next() {
		err := rows.Scan(&p.ID, &p.Text, &p.User_id, &p.User_login)
		if err != nil {
			fmt.Println(err)
			continue
		}
		posts = append(posts, p)
	}
	postjs, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(postjs))
	json.NewEncoder(resp).Encode(posts)
}
func PreAddPost(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Pre-Request")
	resp.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	resp.Header().Set("Access-Control-Allow-Headers", "Content-type")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
}
func AddPost(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("Add Post")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	var new_post Post
	_ = json.NewDecoder(req.Body).Decode(&new_post)
	//fmt.Println(new_post)
	rows, err := db.Query("SELECT id FROM blog_user WHERE login = '" + new_post.User_login + "'")
	rows.Next()
	var u_id int
	rows.Scan(&u_id)
	rows.Close()
	fmt.Println(u_id)
	result, err := db.Exec("insert into post(text, user_id) values (\"" + new_post.Text + "\"," + strconv.Itoa(u_id) + ")")
	fmt.Println("insert into post(text, user_id) values (\"" + new_post.Text + "\"," + strconv.Itoa(u_id) + ")")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
func UserLogin(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Set("Content-Type", "application/json")
	fmt.Println("login")
	var luser User
	_ = json.NewDecoder(req.Body).Decode(&luser)
	fmt.Println(luser)
	var status LoginStatus
	var usr User
	rows, err := db.Query("SELECT login, psw FROM blog_user WHERE login='" + luser.Login + "'")
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		err = rows.Scan(&usr.Login, &usr.Psw)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(usr.Psw)
		if usr.Psw == luser.Psw {
			status.Status = 1
			status.Login = usr.Login
		} else {
			status.Status = 0
		}
	}
	fmt.Println(status)
	json.NewEncoder(resp).Encode(status)
}
func main() {
	db, err = sql.Open("sqlite3", "MyBLog")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/posts/{login}", GetUserPosts).Methods("GET")
	r.HandleFunc("/add_post", AddPost).Methods("POST")
	r.HandleFunc("/add_post", PreAddPost).Methods("OPTIONS")
	r.HandleFunc("/login", PreAddPost).Methods("OPTIONS")
	r.HandleFunc("/login", UserLogin).Methods("POST")
	log.Fatal(http.ListenAndServe(":8101", r))
}
