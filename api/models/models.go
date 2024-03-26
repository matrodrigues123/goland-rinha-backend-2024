package models

import (
	"time"
)

// type Balance struct {
//     Limit int `json:"limite"`
//     Total  int `json:"total"`
//     UpdatedAt time.Time `json:"data_extrato"`
// }

// type AccountResponse struct {
//     Balance Balance `json:"saldo"`
//     Transactions []Transaction `json:"ultimas_transacoes"`
// }

type Account struct {
	ID 		     uint	        
	Limit        int                
	Total        int
}

type Transaction struct {
	ID          uint `json:"-"` 
	AccountID   uint  `json:"-"`  
	Value       int    `json:"valor"` 
	Type        string `json:"tipo"` 
	Description string `json:"descricao"` 
	CreatedAt   time.Time `json:"realizada_em"` 
}

type TransactionRequest struct {
    Value int    `json:"valor" binding:"required,min=1"`
	Type string `json:"tipo" binding:"required,max=1,oneof=d c"`
	Description string `json:"descricao" binding:"required,max=10"`
}
