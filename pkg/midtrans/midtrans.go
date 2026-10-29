package midtrans

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rescuebyte/rescuebyte-be/domain/dto"
)

type Config struct {
	ServerKey    string
	ClientKey    string
	IsProduction bool
}

type Service interface {
	CreateSnapTransaction(ctx context.Context, req dto.SnapTransactionRequest) (*dto.SnapTransactionResponse, error)
	VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool
}

type midtransService struct {
	config     Config
	httpClient *http.Client
}

func NewMidtransService(config Config) Service {
	return &midtransService{
		config: config,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *midtransService) getSnapBaseURL() string {
	if s.config.IsProduction {
		return "https://app.midtrans.com/snap/v1"
	}
	return "https://app.sandbox.midtrans.com/snap/v1"
}

func (s *midtransService) CreateSnapTransaction(ctx context.Context, req dto.SnapTransactionRequest) (*dto.SnapTransactionResponse, error) {
	url := fmt.Sprintf("%s/transactions", s.getSnapBaseURL())

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal snap request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(s.config.ServerKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans snap request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read midtrans snap response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("midtrans api error (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var snapResp dto.SnapTransactionResponse
	if err := json.Unmarshal(respBytes, &snapResp); err != nil {
		return nil, fmt.Errorf("failed to parse midtrans snap response: %w", err)
	}

	return &snapResp, nil
}

func (s *midtransService) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	payload := orderID + statusCode + grossAmount + s.config.ServerKey
	hash := sha512.Sum512([]byte(payload))
	computedSignature := hex.EncodeToString(hash[:])
	return computedSignature == signatureKey
}
