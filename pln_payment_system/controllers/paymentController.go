package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"pln_payment_system/models"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type PaymentController struct {
    db    *sql.DB
    store *sessions.CookieStore
}

func NewPaymentController(db *sql.DB, store *sessions.CookieStore) *PaymentController {
    return &PaymentController{db: db, store: store}
}

func (pc *PaymentController) Index(w http.ResponseWriter, r *http.Request) {
    session, _ := pc.store.Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    rows, err := pc.db.Query("SELECT * FROM payments WHERE user_id = ?", userID)
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

func (pc *PaymentController) Create(w http.ResponseWriter, r *http.Request) {
    session, _ := pc.store.Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var paymentType struct {
        PaymentType string `json:"paymentType"`
    }
    if err := json.NewDecoder(r.Body).Decode(&paymentType); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if paymentType.PaymentType == "monthly" {
        var payment models.Payment
        if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        _, err := pc.db.Exec("INSERT INTO payments (user_id, amount, method, transaction_id) VALUES (?, ?, ?, ?)",
            userID, payment.Amount, payment.Method, payment.TransactionID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    } else if paymentType.PaymentType == "token" {
        var tokenPayment models.TokenPayment
        if err := json.NewDecoder(r.Body).Decode(&tokenPayment); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        _, err := pc.db.Exec("INSERT INTO token_payments (user_id, denomination) VALUES (?, ?)",
            userID, tokenPayment.Denomination)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
    } else {
        http.Error(w, "Invalid payment type", http.StatusBadRequest)
    }
}

func (pc *PaymentController) Update(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid payment ID", http.StatusBadRequest)
        return
    }

    session, _ := pc.store.Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var payment models.Payment
    if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = pc.db.Exec("UPDATE payments SET amount=?, method=?, transaction_id=? WHERE id=? AND user_id=?",
        payment.Amount, payment.Method, payment.TransactionID, id, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (pc *PaymentController) Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid payment ID", http.StatusBadRequest)
        return
    }

    session, _ := pc.store.Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    _, err = pc.db.Exec("DELETE FROM payments WHERE id=? AND user_id=?", id, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
