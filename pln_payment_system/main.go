package main

import (
	"database/sql"
	"log"
	"net/http"
	"pln_payment_system/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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

    router.HandleFunc("/register", userController.Register).Methods("POST")
    router.HandleFunc("/login", userController.Login).Methods("POST")
    router.HandleFunc("/logout", userController.Logout).Methods("POST")
    router.HandleFunc("/user-info", userController.UserInfo).Methods("GET")

    router.HandleFunc("/payments", paymentController.Index).Methods("GET")
    router.HandleFunc("/payments", paymentController.Create).Methods("POST")
    router.HandleFunc("/payments/{id}", paymentController.Update).Methods("PUT")
    router.HandleFunc("/payments/{id}", paymentController.Delete).Methods("DELETE")

    // Static files
    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./views/")))

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
