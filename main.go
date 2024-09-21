package main

import (
    "log"
    "net/http"

    "user_management.com/app/controllers"
    "user_management.com/app/services"
    "database/sql"
    "github.com/rs/cors"
     
)

func main() {
    corsOpts := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, //you service is available and allowed for this base url 
        AllowedMethods: []string{
            http.MethodGet,//http methods for your app
            http.MethodPost,
            http.MethodPut,
            http.MethodPatch,
            http.MethodDelete,
            http.MethodOptions,
            http.MethodHead,
        },
    
        AllowedHeaders: []string{
            "*",//or you can your header key values which you are using in your application
    
        },
    })
    db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/momo_management")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
     // Initialize the user service with the database connection
     services.InitUserService(db)
    router := controllers.SetupRoutes()
    log.Println("Server running on port 9090")
    log.Fatal(http.ListenAndServe(":9090", corsOpts.Handler(router)))
}
