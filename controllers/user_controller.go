package controllers

import (
    "encoding/json"
    "net/http"
    "log"
    "strconv"

    "user_management.com/app/models"
    "user_management.com/app/services"
    "github.com/gorilla/mux"
    "user_management.com/app/middlewares"
    
    
)

func SetupRoutes() *mux.Router {
    router := mux.NewRouter()
    // Public routes
    router.HandleFunc("/user/login", Login).Methods(http.MethodPost)
    router.HandleFunc("/user/sign-up", usersHandler).Methods(http.MethodPost)

    // Protected routes
    protected := router.PathPrefix("/users").Subrouter()
    protected.Use(middlewares.TokenAuthMiddleware)
    protected.HandleFunc("", usersHandler).Methods(http.MethodGet, http.MethodPost)
    protected.HandleFunc("/{userId}", userHandler).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
    return router
}

// func enableCors(w *http.ResponseWriter) {
//     (*w).Header().Set("Access-Control-Allow-Origin", "*")
//     }

func usersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    switch r.Method {
    case http.MethodGet:
        getAllUsers(w)
    case http.MethodPost:
        createUser(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    userId, err := strconv.Atoi(mux.Vars(r)["userId"])
    if err != nil {
        log.Fatal("there is an error with user Id ",err)
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodGet:
        getUserByID(w, userId)
    case http.MethodPut:
        updateUser(w, r, userId)
    case http.MethodDelete:
        deleteUser(w, userId)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getAllUsers(w http.ResponseWriter) {
    apiResponseSuccess:=models.GetUsersResponseSuccess{}
    apiResponseFailure:=models.GetUsersResponseFailure{}
    users,err := services.GetAllUsers()
    if err!=nil{
        log.Println("there is an error with the db ",err)
        apiResponseFailure.ResponseCode="115"
        apiResponseFailure.ResponseMessage="There is a db related issue"
        apiResponseFailure.Meta=err.Error()
        jsonResponse(w, apiResponseFailure, http.StatusInternalServerError)
    }else{
        // b, _ := json.Marshal(users)
        // log.Println("user list ",string(b))
        apiResponseSuccess.ResponseCode="01"
        apiResponseSuccess.ResponseMessage="Successful retrieved Users"
        apiResponseSuccess.Meta=users

    jsonResponse(w, apiResponseSuccess, http.StatusOK)
    }
}

func getUserByID(w http.ResponseWriter, id int) {
    user, err := services.GetUserByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    jsonResponse(w, user, http.StatusOK)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    apiCreateUserSuccessResponse:=models.CreateUserResponseSuccess{}
    apiResponse:=models.Response{}
    log.Println("request comming for create user ",r.Body)
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    log.Println("user found from request ",user)

    newUser,err := services.CreateUser(user)
    if err!=nil{
        log.Println("there is an error with the db ",err)
        apiResponse.ResponseCode="112"
        apiResponse.ResponseMessage="Internal Db Error"
        apiResponse.Meta=err.Error()
        jsonResponse(w, apiResponse, http.StatusInternalServerError)
    }else{
        // b, _ := json.Marshal(newUser)
        apiCreateUserSuccessResponse.ResponseCode="01"
        apiCreateUserSuccessResponse.ResponseMessage="Successful created User"
        apiCreateUserSuccessResponse.Meta=newUser
        jsonResponse(w, apiCreateUserSuccessResponse, http.StatusCreated)
    }
    
}

func updateUser(w http.ResponseWriter, r *http.Request, userId int) {
    var updatedUser models.User
    apiResponse:=models.Response{}
    if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
        // http.Error(w, "Bad request", http.StatusBadRequest)
        apiResponse.ResponseCode="113"
        apiResponse.ResponseMessage="Bad Request Data"
        apiResponse.Meta=err.Error()
        jsonResponse(w, apiResponse, http.StatusBadRequest)
        return
    }
    user, err := services.UpdateUser(userId, updatedUser)
    if err != nil {
        // http.Error(w, err.Error(), http.StatusNotFound)
        apiResponse.ResponseCode="115"
        apiResponse.ResponseMessage="There is a db related issue"
        apiResponse.Meta=err.Error()
        jsonResponse(w, apiResponse, http.StatusNotFound)
        return
    }
    b, _ := json.Marshal(user)
        apiResponse.ResponseCode="01"
        apiResponse.ResponseMessage="Successful updated User"
        apiResponse.Meta=string(b)
    jsonResponse(w, apiResponse, http.StatusOK)
}

func deleteUser(w http.ResponseWriter, userId int) {
    err := services.DeleteUser(userId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
