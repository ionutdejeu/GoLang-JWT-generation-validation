package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/refresh", refreshHandler)

	err := http.ListenAndServe(":9090", nil)
	checkError(err)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signin handler")
	var cr creds
	err := json.NewDecoder(r.Body).Decode(&cr)
	fmt.Println("creds", err, cr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPass, ok := users[cr.Username]

	// when the user exist then we check the passwords
	// else we return 401 Unathorized HTTP
	if !ok || expectedPass != cr.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	cl := &claims{
		Username: cr.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome handler")
	c, er := r.Cookie("token")
	if er != nil {
		if er == http.ErrNoCookie {
			// unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookieValue := c.Value
	cl := &claims{}

	// parse the jwt string
	tkn, err := jwt.ParseWithClaims(cookieValue, cl, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Successfull login %s!", cl.Username)))
}
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Refres handler")
}
func checkError(err error) {
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
