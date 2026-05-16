package handlers

import (
	"net/http"
	"strconv"
	"wallet-api/store"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

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

	if err := store.SaveWallet(*wallet); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save wallet"})
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

	wallet, err := store.GetWallet(id)
	if err != nil {
		return c.JSON(404, map[string]string{
			"error": "wallet not found",
		})
	}

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

	wallet, err := store.GetWallet(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Wallet Not Found",
		})
	}

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

	if err := store.SaveWallet(*wallet); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save transaction"})
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

	wallet, err := store.GetWallet(id)
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
	start := offset
	end := offset + limit

	if start > len(wallet.Transactions) {
		start = len(wallet.Transactions)
	}

	if end > len(wallet.Transactions) {
		end = len(wallet.Transactions)
	}

	return c.JSON(http.StatusOK, wallet.Transactions[start:end])
}
