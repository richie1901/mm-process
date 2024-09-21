package models




// Struct to hold the login request body
type Credentials struct {
    Email string `json:"email"`
    Password string `json:"password"`
}
