package midtrans_test

import (
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rescuebyte/rescuebyte-be/pkg/midtrans"
)

func TestVerifySignature(t *testing.T) {
	serverKey := "SB-Mid-server-test12345"
	svc := midtrans.NewMidtransService(midtrans.Config{
		ServerKey:    serverKey,
		ClientKey:    "SB-Mid-client-test12345",
		IsProduction: false,
	})

	orderID := "RB251026-00123"
	statusCode := "200"
	grossAmount := "22000.00"

	payload := orderID + statusCode + grossAmount + serverKey
	hash := sha512.Sum512([]byte(payload))
	validSignature := hex.EncodeToString(hash[:])

	// Should pass with correct signature
	assert.True(t, svc.VerifySignature(orderID, statusCode, grossAmount, validSignature))

	// Should fail with modified signature
	assert.False(t, svc.VerifySignature(orderID, statusCode, grossAmount, "invalid-signature"))

	// Should fail with altered amount
	assert.False(t, svc.VerifySignature(orderID, statusCode, "50000.00", validSignature))
}
