package models

type Payment struct {
    ID            int     `json:"id"`
    UserID        int     `json:"user_id"`
    Amount        float64 `json:"amount"`
    Method        string  `json:"method"`
    TransactionID string  `json:"transaction_id"`
}

type TokenPayment struct {
    ID           int `json:"id"`
    UserID       int `json:"user_id"`
    Denomination int `json:"denomination"`
}
