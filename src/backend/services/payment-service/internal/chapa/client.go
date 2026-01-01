package chapa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	ChapaTestBaseURL = "https://api.chapa.co/v1"
	ChapaProdBaseURL = "https://api.chapa.co/v1"
)

type Client struct {
	SecretKey     string
	PublicKey     string
	EncryptionKey string
	BaseURL       string
	HTTPClient    *http.Client
}

type InitializePaymentRequest struct {
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	Email       string            `json:"email"`
	FirstName   string            `json:"first_name"`
	LastName    string            `json:"last_name"`
	TxRef       string            `json:"tx_ref"`
	CallbackURL string            `json:"callback_url"`
	ReturnURL   string            `json:"return_url"`
	CustomTitle string            `json:"customization[title],omitempty"`
	CustomDesc  string            `json:"customization[description],omitempty"`
	Metadata    map[string]string `json:"meta,omitempty"`
}

type InitializePaymentResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

type VerifyPaymentResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		Amount     float64 `json:"amount"`
		Currency   string  `json:"currency"`
		Email      string  `json:"email"`
		FirstName  string  `json:"first_name"`
		LastName   string  `json:"last_name"`
		TxRef      string  `json:"tx_ref"`
		Status     string  `json:"status"`
		Reference  string  `json:"reference"`
		ChargeType string  `json:"charge_type"`
		CreatedAt  string  `json:"created_at"`
		UpdatedAt  string  `json:"updated_at"`
	} `json:"data"`
}

func NewClient(secretKey, publicKey, encryptionKey string, isTest bool) *Client {
	baseURL := ChapaProdBaseURL
	if isTest {
		baseURL = ChapaTestBaseURL
	}

	return &Client{
		SecretKey:     secretKey,
		PublicKey:     publicKey,
		EncryptionKey: encryptionKey,
		BaseURL:       baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) InitializePayment(req InitializePaymentRequest) (*InitializePaymentResponse, error) {
	// Generate unique transaction reference if not provided
	if req.TxRef == "" {
		req.TxRef = fmt.Sprintf("RF-%s", uuid.New().String()[:8])
	}

	// Set default currency
	if req.Currency == "" {
		req.Currency = "ETB"
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/transaction/initialize", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.SecretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chapa API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chapaResp InitializePaymentResponse
	if err := json.Unmarshal(body, &chapaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chapaResp.Status != "success" {
		return nil, fmt.Errorf("chapa initialization failed: %s", chapaResp.Message)
	}

	return &chapaResp, nil
}

func (c *Client) VerifyPayment(txRef string) (*VerifyPaymentResponse, error) {
	httpReq, err := http.NewRequest("GET", c.BaseURL+"/transaction/verify/"+txRef, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.SecretKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chapa API error (status %d): %s", resp.StatusCode, string(body))
	}

	var chapaResp VerifyPaymentResponse
	if err := json.Unmarshal(body, &chapaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &chapaResp, nil
}
