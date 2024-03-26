package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/matrodrigues123/rinha-2024-go-api/database"
	"github.com/matrodrigues123/rinha-2024-go-api/models"
)


func main() {
	// defininig time zone
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	time.Local = loc

	// postgres db
	db_postgres := database.Connection()
	
	r := gin.Default()
  
  	r.POST("/clientes/:id/transacoes", func (c *gin.Context) {
		postTransaction(c, db_postgres)
	})
  	r.GET("/clientes/:id/extrato", func (c *gin.Context) {
		getBalance(c, db_postgres)
	})
  
  	r.Run(":3000")
    
}

func postTransaction(c *gin.Context, db *gorm.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid ID. Must be an integer."})
		return
	}

	var newTransactionRequest models.TransactionRequest
	if err := c.BindJSON(&newTransactionRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
        return
    }
	
	tx := db.Begin()
	
	var account models.Account
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found."})
		return
	}

	switch newTransactionRequest.Type {
	case "d":
		newTotal := account.Total - newTransactionRequest.Value
		if (newTotal + account.Limit < 0) {
			tx.Rollback()
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid debit operation"})
			return
		}
		account.Total = newTotal
	case "c":
		account.Total += newTransactionRequest.Value
	default:
		tx.Rollback()
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Operation must be c or d"})
		return
	}

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account balance."})
		return
	}

	newTransaction := models.Transaction{
		AccountID: account.ID,
		Value: newTransactionRequest.Value,
		Type: newTransactionRequest.Type,
		Description: newTransactionRequest.Description,
	}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Create(&newTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record new transaction."})
		return
	}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed."})
		return
	}

	response := gin.H{
		"limite": account.Limit,
		"saldo":  account.Total,
	}

	c.JSON(http.StatusOK, response)
}

func getBalance(c *gin.Context, db *gorm.DB)  {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID. Must be an integer."})
		return
	}

	var account models.Account
    
	if err := db.Clauses(clause.Locking{Strength: "SHARE"}).First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID or Account not found."})
		return
	}

	var transactions []models.Transaction
	if err := db.Clauses(clause.Locking{Strength: "SHARE"}).
		Select("value, type, description, created_at").
		Where("account_id = ?", account.ID).
		Order("created_at desc").
		Limit(10).
		Find(&transactions).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldnt find transactions"})
			return
		}
	

	response := gin.H{
        "saldo": gin.H{
            "limite":     account.Limit,
            "total":      account.Total,
            "data_extrato": time.Now(),
        },
        "ultimas_transacoes": transactions,
    }

    c.JSON(http.StatusOK, response)

}
