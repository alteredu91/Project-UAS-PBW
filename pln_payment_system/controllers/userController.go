package controllers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "golang.org/x/crypto/bcrypt"
    "github.com/gorilla/sessions"
)

type UserController struct {
    db    *sql.DB
    store *sessions.CookieStore
}

func NewUserController(db *sql.DB, store *sessions.CookieStore) *UserController {
    return &UserController{db: db, store: store}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
    var user struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    _, err = uc.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
    var user struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var hashedPassword string
    var userID int
    err := uc.db.QueryRow("SELECT id, password FROM users WHERE username = ?", user.Username).Scan(&userID, &hashedPassword)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    session, _ := uc.store.Get(r, "session")
    session.Values["user_id"] = userID
    session.Save(r, w)

    w.WriteHeader(http.StatusOK)
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := uc.store.Get(r, "session")
    delete(session.Values, "user_id")
    session.Save(r, w)

    w.WriteHeader(http.StatusOK)
}

func (uc *UserController) UserInfo(w http.ResponseWriter, r *http.Request) {
    session, _ := uc.store.Get(r, "session")
    userID, ok := session.Values["user_id"].(int)
    if !ok {
        json.NewEncoder(w).Encode(map[string]interface{}{"loggedIn": false})
        return
    }

    var username string
    err := uc.db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"loggedIn": true, "username": username})
}

func (uc *UserController) AdminLogin(w http.ResponseWriter, r *http.Request) {
    var admin struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Hardcoded admin credentials
    const hardcodedAdminUsername = "admin"
    const hardcodedAdminPassword = "admin"

    if admin.Username != hardcodedAdminUsername || admin.Password != hardcodedAdminPassword {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    session, _ := uc.store.Get(r, "admin_session")
    session.Values["admin_id"] = hardcodedAdminUsername
    session.Save(r, w)

    w.WriteHeader(http.StatusOK)
}
