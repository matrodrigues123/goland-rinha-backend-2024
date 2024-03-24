package models

import (
	"sync"
	"time"
)

type Balance struct {
    Limit int `json:"limite"`
    Total  int `json:"total"`
    UpdatedAt time.Time `json:"data_extrato"`
}

type Account struct {
    Mutex sync.RWMutex `json:"-"`
    Balance Balance `json:"saldo"`
    Transactions []Transaction `json:"ultimas_transacoes"`
}

type Transaction struct {
    Value int    `json:"valor" binding:"required,min=1"`
	Type string `json:"tipo" binding:"required,max=1,oneof=d c"`
	Description string `json:"descricao" binding:"required,max=10"`
    CreatedAt time.Time `json:"realizada_em"`
}
