package controllers

import (
    "encoding/json"
    "net/http"
    "time"
	"log"
    "github.com/golang-jwt/jwt"
    "user_management.com/app/models"
    "user_management.com/app/services"
)

var jwtSecret = []byte("secretManagement")  // Change this to a secure key


// Login handler to authenticate user and return a JWT token
func Login(w http.ResponseWriter, r *http.Request) {
    var creds models.Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    // Validate the user credentials
    user, err := services.AuthenticateUser(creds.Email, creds.Password)
    if err != nil {
		log.Fatal("error found authenticating ",err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Generate JWT token
    token, err := generateJWT(user)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    jsonResponse(w, map[string]string{"token": token}, http.StatusOK)
}

// Function to generate JWT token
func generateJWT(user models.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": user.Email,
        "exp":      time.Now().Add(time.Hour * 2).Unix(),  // Token expiry time (2 hours)
    })

    // Sign the token with the secret key
    return token.SignedString(jwtSecret)
}
