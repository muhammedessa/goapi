package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	Name string
	Email string
}

func initialMigration()  {
	db,err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Database failed to connect")
	}
	defer db.Close()
	db.AutoMigrate(&User{})
}

func handleRequests(){
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/users" , allUsers).Methods("GET")
	myRouter.HandleFunc("/user/{id}" , showUser).Methods("GET")
	myRouter.HandleFunc("/user/{id}" , deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/user/{id}" , updateUser).Methods("PUT")
	myRouter.HandleFunc("/user" , newUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081",myRouter))
}

func  main()  {
	fmt.Println("Go CRUD API with SQLite3")
	initialMigration()
	handleRequests()
}

func allUsers(w http.ResponseWriter , r *http.Request)  {

	db,err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Database failed to connect")
	}
	defer db.Close()
	var users []User
	db.Find(&users)
	fmt.Println(users)
	json.NewEncoder(w).Encode(users)
}
func newUser(w http.ResponseWriter , r *http.Request)  {
	db,err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Database failed to connect")
	}
	defer db.Close()
	var input User
	errs := json.NewDecoder(r.Body).Decode(&input)
	if errs != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w,"Please check your json format")
		return
	}
	db.Create(&User{Name: input.Name, Email: input.Email})
	fmt.Println("New record added successfully !")
}

func deleteUser(w http.ResponseWriter , r *http.Request)  {
	db,err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Database failed to connect")
	}
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	var user User
	db.Where("ID =?", id).Find(&user)
	db.Delete(&user)
	fmt.Fprintf(w,"Record deleted successfully !")
}
func updateUser(w http.ResponseWriter , r *http.Request)  {
	db,err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Database failed to connect")
	}
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	var user User
	db.Where("ID =?", id).Find(&user)
	var input User
	errs := json.NewDecoder(r.Body).Decode(&input)
	if errs != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w,"Please check your json format")
		return
	}
	user.Name = input.Name
	user.Email = input.Email
	db.Save(&user)
	fmt.Fprintf(w,"Record updated successfully !")
}
func showUser(w http.ResponseWriter , r *http.Request)  {
	db,err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("Database failed to connect")
	}
	defer db.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	var user User
	db.Where("ID =?", id).Find(&user)

	json.NewEncoder(w).Encode(user)
}
