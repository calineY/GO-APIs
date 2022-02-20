package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)


type User struct{
	Name string `json:"name"`
	Email string `json:"user_email"`
	Password string `jason:"password"`
}

var db *sql.DB
var err error

func main(){
	fmt.Println("GO")
	db, err= sql.Open("mysql", "root:@/facebookdb")

	if err!=nil{
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("connected")

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/signup",signUp).Methods("POST")
 	http.ListenAndServe(":8000", router)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var users []User
	
	results, err := db.Query("SELECT user_name,user_email,password FROM users")

	if err!=nil{
		panic(err.Error())
	}

	defer results.Close()

	for results.Next(){
		var user User

		err=results.Scan(&user.Name,&user.Email,&user.Password)

		if err != nil{
			panic(err.Error())
		}
		users = append(users, User{Name:user.Name, Email: user.Email, Password: user.Password})

	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT user_name,user_email,password FROM users WHERE user_id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var user User
	for result.Next() {
		err := result.Scan(&user.Name,&user.Email,&user.Password)
		if err != nil {
		panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(user)
}

func signUp(w http.ResponseWriter, r *http.Request) {  
	
	stmt, err := db.Prepare("INSERT INTO users(user_name,user_email,password) VALUES (?,?,?)")
	  if err != nil {
		panic(err.Error())
	  }  
	  body, err := ioutil.ReadAll(r.Body)
	  if err != nil {
		panic(err.Error())
	  }
	  body_str:=string(body)
	  fmt.Println(body_str)
	  keyVal := make(map[string]string)
	  json.Unmarshal(body, &keyVal)
	  user_name := keyVal["user_name"]
	  user_email := keyVal["user_email"]
	  password := keyVal["password"]
	  
	  _, err = stmt.Exec(user_name, user_email, password)
	
	  if err != nil {
		panic(err.Error())
	  }  
	  fmt.Fprintf(w, "User added")
} 
