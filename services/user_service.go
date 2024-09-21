package services

import (
    "errors"
    "log"
    "sync"
    "database/sql"

    "user_management.com/app/models"
    _ "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"
)

var (
    db *sql.DB
    mu    sync.Mutex  // this is used to handle concurrent requests
)

// var db *sql.DB

// Initialize the user service with the database connection
func InitUserService(database *sql.DB) {
    db = database
}

// Fetch all users from the database
func GetAllUsers() ([]models.User, error) {
    rows, err := db.Query("SELECT * FROM user")
    if err != nil {
        log.Println("there is an error fetching from db ",err)
        return nil, err
    }
    defer rows.Close()
    var users []models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.UserId,&user.FullName,&user.DateOfBirth,&user.Address,&user.MobileNumber,&user.Email,&user.DateCreated,&user.Password); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}

// Insert a new user into the database
func CreateUser(user models.User) (models.User, error) {

       // Hash the password
       hashedPassword, err := hashPassword(user.Password)
       if err != nil {
           return models.User{}, err
       }

    result, err := db.Exec("INSERT INTO user (fullname, date_of_birth,address,mobile_number,email,password) VALUES (?, ?,?,?,?,?)", user.FullName, user.DateOfBirth,user.Address,user.MobileNumber,user.Email,hashedPassword)
    
    if err != nil {
        return models.User{}, err
    }

    userID, err := result.LastInsertId()
    if err != nil {
        return models.User{}, err
    }
    user.UserId = int(userID)
    user.Password=hashedPassword
    return user, nil
}

// Fetch a user by userId
func GetUserByID(userId int) (models.User, error) {
    var user models.User
    err := db.QueryRow("SELECT user_id, fullname, date_of_birth,address,mobile_number,email,date_created,password FROM user WHERE user_id = ?", userId).Scan(&user.UserId, &user.FullName, &user.DateOfBirth,&user.Address,&user.MobileNumber,&user.Email,&user.DateCreated,&user.Password)
    if err == sql.ErrNoRows {
        return models.User{}, errors.New("user not found")
    } else if err != nil {
        return models.User{}, err
    }
    return user, nil
}

// Update a user's details in the database
func UpdateUser(userId int, updatedUser models.User) (models.User, error) {
    result, err := db.Exec("UPDATE user SET fullname = ?, date_of_birth = ?,address = ?,mobile_number = ? WHERE user_id = ?", updatedUser.FullName, updatedUser.DateOfBirth,updatedUser.Address,updatedUser.MobileNumber, userId)
    if err != nil {
        log.Println("there is an error in db",err)
        return models.User{}, err
    }
     // Check how many rows were affected
     rowsAffected, err := result.RowsAffected()
     if err != nil {
        log.Println("there is an error in db rows",err)
        return models.User{}, err
     }
 
     if rowsAffected == 0 {
        return models.User{}, errors.New("No Update Effected")
     }
     // Hash the password
     hashedPassword, err := hashPassword(updatedUser.Password)
     if err != nil {
         return models.User{}, err
     }
    updatedUser.UserId = userId
    updatedUser.Password=hashedPassword
    return updatedUser, nil
}

// Delete a user from the database
func DeleteUser(userId int) error {
    _, err := db.Exec("DELETE FROM user WHERE user_id = ?", userId)
    if err != nil {
        return err
    }
    return nil
}

// Authenticate a user by username and password
func AuthenticateUser(username, password string) (models.User, error) {
    var user models.User
    var hashedPassword string

    // Query the user and the hashed password from the database
    err := db.QueryRow("SELECT user_id, fullname, date_of_birth,address,mobile_number,email,date_created,password FROM user WHERE email = ?", username).Scan(&user.UserId, &user.FullName, &user.DateOfBirth,&user.Address,&user.MobileNumber,&user.Email,&user.DateCreated,&hashedPassword)
    if err == sql.ErrNoRows {
        return models.User{}, errors.New("invalid username or password")
    } else if err != nil {
        return models.User{}, err
    }

    // Compare the hashed password stored in the database with the password provided by the user
    err = verifyPassword(hashedPassword, password)
    if err != nil {
        return models.User{}, errors.New("invalid username or password")
    }

    return user, nil
}

// Hash the password before storing it in the database
func hashPassword(password string) (string, error) {
    // bcrypt.GenerateFromPassword generates a hashed password
    if password == "" || password == "null"{
        return "",errors.New("Password cannot be null or empty")
    }
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// VerifyPassword compares a hashed password with the provided password
func verifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
