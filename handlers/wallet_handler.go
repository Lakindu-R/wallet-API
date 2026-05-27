package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"wallet-api/config"
	"wallet-api/models"
)

// CreateWallet godoc
// @Summary Create a new wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param input body object true "Wallet name" SchemaExample({"name":"My Wallet"})
// @Success 200 {object} models.Wallet
// @Failure 400 {object} map[string]string
// @Router /wallets [post]
func CreateWallet(c echo.Context) error {
	var input struct {
		Name string `json:"name"`
	}

	if err := c.Bind(&input); err != nil || input.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	id := uuid.New().String()
	wallet := &models.Wallet{
		ID:           id,
		Name:         input.Name,
		Balance:      0,
		Transactions: []models.Transaction{},
	}

	walletJSON, _ := json.Marshal(wallet)

	err := config.RedisClient.Set(
		config.Ctx,
		"wallet:"+wallet.ID,
		walletJSON,
		0,
	).Err()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save wallet",
		})
	}
	return c.JSON(http.StatusOK, wallet)
}

// GetWallet godoc
// @Summary Get a wallet by ID
// @Tags wallets
// @Produce json
// @Param id path string true "Wallet ID"
// @Success 200 {object} models.Wallet
// @Failure 404 {object} map[string]string
// @Router /wallets/{id} [get]
func GetWallet(c echo.Context) error {
	id := c.Param("id")

	result, err := config.RedisClient.Get(
		config.Ctx,
		"wallet:"+id,
	).Result()

	if err != nil {
		return err
	}

	var wallet models.Wallet

	json.Unmarshal([]byte(result), &wallet)

	return c.JSON(http.StatusOK, wallet)
}

// AddTransaction godoc
// @Summary Add a transaction to a wallet
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Wallet ID"
// @Param input body object true "Transaction input" SchemaExample({"type":"credit","amount":100})
// @Success 200 {object} models.Wallet
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /wallets/{id}/transactions [post]
func AddTransaction(c echo.Context) error {
	id := c.Param("id")

	result, err := config.RedisClient.Get(
		config.Ctx,
		"wallet:"+id,
	).Result()

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Wallet Not Found",
		})
	}

	var wallet models.Wallet
	json.Unmarshal([]byte(result), &wallet)

 	var input struct {
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
	}

	if err := c.Bind(&input); err != nil || input.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid Input",
		})
	}

 	if input.Type == "credit" {
		wallet.Balance += input.Amount
	} else if input.Type == "debit" {
		if wallet.Balance < input.Amount {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Insufficient funds",
			})
		}
		wallet.Balance -= input.Amount
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid transaction type",
		})
	}

 	tx := models.Transaction{
		ID:     uuid.New().String(),
		Type:   input.Type,
		Amount: input.Amount,
	}

	wallet.Transactions = append(wallet.Transactions, tx)

	txJSON, _ := json.Marshal(tx)

	config.RedisClient.RPush(
		config.Ctx,
		"transactions:"+id,
		txJSON,
	)

	walletJSON, _ := json.Marshal(wallet)

	err = config.RedisClient.Set(
		config.Ctx,
		"wallet:"+id,
		walletJSON,
		0,
	).Err()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save transaction",
		})
	}

	return c.JSON(http.StatusOK, wallet)
}

// GetTransactions godoc
// @Summary Get transactions for a wallet
// @Tags transactions
// @Produce json
// @Param id path string true "Wallet ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Transaction
// @Failure 404 {object} map[string]string
// @Router /wallets/{id}/transactions [get]
func GetTransactions(c echo.Context) error {
	id := c.Param("id")

	// check wallet exists (optional but good)
	_, err := config.RedisClient.Get(
		config.Ctx,
		"wallet:"+id,
	).Result()

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Wallet not found",
		})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 5
	}

	// 🔥 Redis pagination using LRange
	result, err := config.RedisClient.LRange(
		config.Ctx,
		"transactions:"+id,
		int64(offset),
		int64(offset+limit-1),
	).Result()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch transactions",
		})
	}

	// convert JSON strings → structs
	var transactions []models.Transaction

	for _, item := range result {
		var tx models.Transaction
		json.Unmarshal([]byte(item), &tx)
		transactions = append(transactions, tx)
	}

	return c.JSON(http.StatusOK, transactions)
}
