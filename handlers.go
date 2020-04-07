package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_encoding_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type creds struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/refresh", refreshHandler)

	err := http.ListenAndServe(":9090", nil)
	checkError(err)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signin handler")
}
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome handler")
}
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Refres handler")
}
func checkError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
