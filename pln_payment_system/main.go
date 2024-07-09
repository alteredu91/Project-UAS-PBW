package main

import (
    "database/sql"
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "pln_payment_system/controllers"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func main() {
    router := mux.NewRouter()

    // Database connection
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/pln_payment")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Controllers
    userController := controllers.NewUserController(db, store)
    paymentController := controllers.NewPaymentController(db, store)
    adminController := controllers.NewAdminController(db, store) // Added AdminController

    router.HandleFunc("/register", userController.Register).Methods("POST")
    router.HandleFunc("/login", userController.Login).Methods("POST")
    router.HandleFunc("/logout", userController.Logout).Methods("POST")
    router.HandleFunc("/user-info", userController.UserInfo).Methods("GET")

    router.HandleFunc("/payments", paymentController.Index).Methods("GET")
    router.HandleFunc("/payments", paymentController.Create).Methods("POST")
    router.HandleFunc("/payments/{id}", paymentController.Update).Methods("PUT")
    router.HandleFunc("/payments/{id}", paymentController.Delete).Methods("DELETE")

    // Admin routes
    router.HandleFunc("/admin/users", adminController.IndexUsers).Methods("GET")
    router.HandleFunc("/admin/payments", adminController.IndexPayments).Methods("GET")

    // Static files
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./views/")))

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
