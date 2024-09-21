package models

// import (
//     "time"

//     _ "github.com/go-sql-driver/mysql"
// )

type User struct{
	UserId int `json:"userId"`
	FullName string `json:"fullName"`
	DateOfBirth string `json:"dateOfBirth"`
	Address string `json:"address"`
	MobileNumber string `json:"mobileNumber"`
	Email string `json:"email"`
	Password string `json:"password"`
	DateCreated string `json:"dateCreated"`
}

// CREATE TABLE user (
//     user_id INT AUTO_INCREMENT PRIMARY KEY,
//     fullname VARCHAR(100) NOT NULL,
//     date_of_birth DATE NOT NULL,
//     address VARCHAR(100) NOT NULL,
//     mobile_number VARCHAR(12) NOT NULL,
//     email VARCHAR(100) NOT NULL,
//     date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     UNIQUE INDEX (email),
//     UNIQUE INDEX (mobile_number)
// );