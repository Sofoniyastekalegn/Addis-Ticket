package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChapaPaymentRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	TxRef       string `json:"tx_ref"`
	CallbackURL string `json:"callback_url"`
	ReturnURL   string `json:"return_url"`
}

type ChapaResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
		Reference   string `json:"reference"`
	} `json:"data"`
}

func handleInitiatePayment(c *gin.Context) {
	var req struct {
		Amount   string `json:"amount" binding:"required"`
		Currency string `json:"currency" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		TxRef    string `json:"tx_ref" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	// Prepare Chapa payment request
	chapaReq := ChapaPaymentRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Email:       req.Email,
		FirstName:   userName.(string),
		LastName:    "",
		TxRef:       req.TxRef,
		CallbackURL: "http://localhost:8080/api/payment/verify",
		ReturnURL:   "http://localhost:3000/payment/success",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(chapaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare payment request"})
		return
	}

	// Make request to Chapa API
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.chapa.co/v1/transaction/initialize", strings.NewReader(string(jsonData)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("CHAPA_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment service unavailable"})
		return
	}
	defer resp.Body.Close()

	var chapaResp ChapaResponse
	if err := json.NewDecoder(resp.Body).Decode(&chapaResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment response"})
		return
	}

	if chapaResp.Status != "success" {
		c.JSON(http.StatusBadRequest, gin.H{"error": chapaResp.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkout_url": chapaResp.Data.CheckoutURL,
		"reference":    chapaResp.Data.Reference,
	})
}
