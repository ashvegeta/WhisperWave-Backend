package handlers

import (
	"WhisperWave-BackEnd/models"
	"WhisperWave-BackEnd/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Handler for logging in existing users
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.UserLoginCredentials

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check for user in DB and generate a session token
	if user.UserName == "ashik" && user.Password == "123456" {
		tokenString, err := utils.GenerateToken(user.UserName)
		if err != nil {
			log.Println("Error Generating JWT Token: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
		}
		
		// return the generated token
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tokenString))
		return

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid Credentials")
	}
}

// Handler for signing up new users
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserSignupCredentials

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error : Bad Request"))
		return
	}

	// store the user on database
	// -------- Write the DB logic here ---------- //


	// Generate session token for the user
	token, err := utils.GenerateToken(user.UserName)
	if err != nil {
		log.Println("Error Generating JWT Token: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// return the session token
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

// handler to check validity of token (ONLY MEANT FOR TESTING)
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")

	err := utils.VerifyToken(token)

	if err!=nil{
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("token is not valid"))
		return
	}

	w.Write([]byte("token is valid"))
}