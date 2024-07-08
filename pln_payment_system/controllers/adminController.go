package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "github.com/gorilla/sessions"
    "pln_payment_system/models"
)

type AdminController struct {
    db    *sql.DB
    store *sessions.CookieStore
}

func NewAdminController(db *sql.DB, store *sessions.CookieStore) *AdminController {
    return &AdminController{db: db, store: store}
}

func (ac *AdminController) AdminUserInfo(w http.ResponseWriter, r *http.Request) {
    session, _ := ac.store.Get(r, "admin_session")
    adminID, ok := session.Values["admin_id"].(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var username string
    err := ac.db.QueryRow("SELECT username FROM admins WHERE id = ?", adminID).Scan(&username)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"loggedIn": true, "username": username})
}

func (ac *AdminController) IndexUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := ac.db.Query("SELECT id, username FROM users")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    users := []models.User{}
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Username); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func (ac *AdminController) IndexPayments(w http.ResponseWriter, r *http.Request) {
    rows, err := ac.db.Query("SELECT * FROM payments")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    payments := []models.Payment{}
    for rows.Next() {
        var payment models.Payment
        if err := rows.Scan(&payment.ID, &payment.UserID, &payment.Amount, &payment.Method, &payment.TransactionID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        payments = append(payments, payment)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(payments)
}
