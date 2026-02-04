package payment

import (
	"fmt"

	"github.com/UthitSawatdee/GoMarketAPI/internal/port"
	"github.com/google/uuid"
)

// MockPaymentService is a mock implementation for testing and development
// Replace this with actual Stripe/PromptPay integration in production
type MockPaymentService struct{}

// NewMockPaymentService creates a new mock payment service
func NewMockPaymentService() port.PaymentService {
	return &MockPaymentService{}
}

// ProcessPayment simulates payment processing
// In production, this would call Stripe/PromptPay/etc.
func (s *MockPaymentService) ProcessPayment(req port.PaymentRequest) (*port.PaymentResult, error) {
	// Validate request
	if req.Amount <= 0 {
		return &port.PaymentResult{
			Status:  "failed",
			Message: "Invalid payment amount",
		}, fmt.Errorf("payment amount must be positive")
	}

	if req.Method == "" {
		return &port.PaymentResult{
			Status:  "failed",
			Message: "Payment method required",
		}, fmt.Errorf("payment method is required")
	}

	// Simulate different payment methods
	switch req.Method {
	case "credit_card":
		// Simulate successful credit card payment
		return &port.PaymentResult{
			TransactionID: uuid.New().String(),
			Status:        "success",
			Message:       "Credit card payment processed successfully",
		}, nil

	case "promptpay":
		// Simulate PromptPay pending (requires QR scan)
		return &port.PaymentResult{
			TransactionID: uuid.New().String(),
			Status:        "pending",
			Message:       "Waiting for PromptPay confirmation",
			RedirectURL:   fmt.Sprintf("/payment/promptpay/%d", req.OrderID),
		}, nil

	case "cod":
		// Cash on delivery - always succeeds
		return &port.PaymentResult{
			TransactionID: uuid.New().String(),
			Status:        "success",
			Message:       "Cash on Delivery - payment will be collected upon delivery",
		}, nil

	default:
		return &port.PaymentResult{
			Status:  "failed",
			Message: fmt.Sprintf("Unsupported payment method: %s", req.Method),
		}, fmt.Errorf("unsupported payment method: %s", req.Method)
	}
}

// VerifyWebhook validates webhook signatures from payment providers
func (s *MockPaymentService) VerifyWebhook(payload []byte, signature string) (bool, error) {
	// Mock always returns true for testing
	// In production, verify HMAC signature from Stripe/PromptPay
	return true, nil
}

// GetPaymentStatus retrieves payment status by transaction ID
func (s *MockPaymentService) GetPaymentStatus(transactionID string) (*port.PaymentResult, error) {
	// Mock always returns success
	// In production, query payment provider's API
	return &port.PaymentResult{
		TransactionID: transactionID,
		Status:        "success",
		Message:       "Payment confirmed",
	}, nil
}
