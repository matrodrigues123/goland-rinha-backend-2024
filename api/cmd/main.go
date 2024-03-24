package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/matrodrigues123/rinha-2024-go-api/models"
)



func main() {

	// defininig time zone
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	time.Local = loc

	// in memory db
	db := map[int]*models.Account{
		1: {
			Balance: models.Balance{
				Limit:        1e5,
				Total: 0,
			},
			Transactions: []models.Transaction{},
		},
		2: {
			Balance: models.Balance{
				Limit:        8e4,
				Total: 0,
			},
			Transactions: []models.Transaction{},
		},
		3: {
			Balance: models.Balance{
				Limit:        1e6,
				Total: 0,
			},
			Transactions: []models.Transaction{},
		},
		4: {
			Balance: models.Balance{
				Limit:        1e7,
				Total: 0,
			},
			Transactions: []models.Transaction{},
		},
		5: {
			Balance: models.Balance{
				Limit:        5e5,
				Total: 0,
			},
			Transactions: []models.Transaction{},
		},
	}
	
	r := gin.Default()
  
  	r.POST("/clientes/:id/transacoes", func (c *gin.Context) {
		postTransaction(c, db)
	})
  	r.GET("/clientes/:id/extrato", func (c *gin.Context) {
		getBalance(c, db)
	})
  
  	r.Run(":3000")
    
}

func postTransaction(c *gin.Context, db map[int]*models.Account) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid ID. Must be an integer."})
		return
	}

	var newTransaction models.Transaction
	if err := c.BindJSON(&newTransaction); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
        return
    }
	
	account, ok := db[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	account.Mutex.Lock()
	defer account.Mutex.Unlock()

	switch newTransaction.Type {
	case "d":
		newTotal := account.Balance.Total - newTransaction.Value
		if (newTotal + account.Balance.Limit < 0) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid debit operation"})
			return
		}
		account.Balance.Total = newTotal
	case "c":
		account.Balance.Total += newTransaction.Value
	default:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Operation must be c or d"})
		return
	}

	newTransaction.CreatedAt = time.Now()
	account.Transactions = append([]models.Transaction{newTransaction}, account.Transactions...)

	response := gin.H{
		"limite": account.Balance.Limit,
		"saldo":  account.Balance.Total,
	}

	c.JSON(http.StatusOK, response)

}

func getBalance(c *gin.Context, db map[int]*models.Account)  {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID. Must be an integer."})
		return
	}

	account, ok := db[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
		return
	}
	account.Mutex.RLock()
	defer account.Mutex.RUnlock()

	account.Balance.UpdatedAt = time.Now()
	
	c.JSON(http.StatusOK, account)
}
