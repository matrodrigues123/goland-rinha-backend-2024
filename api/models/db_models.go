package models

import (
	"time"
)

type AccountsDB struct {
	ID 		     uint 		        `json:"-" gorm:"primaryKey"`
	Limit        int                `json:"limite" gorm:"column:limit;not null"`
	Total        int                `json:"total" gorm:"column:total"`
	VisitedAt    time.Time          `json:"data_extrato" gorm:"-"`
	Transactions []TransactionsDB   `json:"ultimas_transacoes" gorm:"foreignKey:AccountID"`
}

type TransactionsDB struct {
	ID          uint      `json:"-" gorm:"primaryKey"` 
	AccountID   uint      `json:"-" gorm:"column:account_id"`
	Value       int       `json:"valor" gorm:"column:valor;not null" binding:"required,min=1"`
	Type        string    `json:"tipo" gorm:"column:type;type:char(1);not null" binding:"required,oneof=d c"`
	Description string    `json:"descricao" gorm:"column:description;type:varchar(10);not null" binding:"required,max=10"`
	CreatedAt   time.Time `json:"realizada_em" gorm:"column:created_at;autoCreateTime"`
}
