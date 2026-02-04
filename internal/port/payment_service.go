package port

// PaymentRequest contains the information needed to process a payment
type PaymentRequest struct {
	OrderID  uint    `json:"order_id"`
	Amount   float64 `json:"amount"`
	Method   string  `json:"method"`   // "credit_card", "promptpay", "cod"
	Currency string  `json:"currency"` // "THB", "USD"
}

// PaymentResult contains the result of a payment processing attempt
type PaymentResult struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"` // "success", "pending", "failed"
	Message       string `json:"message"`
	RedirectURL   string `json:"redirect_url,omitempty"` // For payment gateways requiring redirect
}

// PaymentService defines the interface for payment processing
// This allows swapping between mock, Stripe, PromptPay, etc.
type PaymentService interface {
	// ProcessPayment initiates a payment for the given request
	ProcessPayment(req PaymentRequest) (*PaymentResult, error)

	// VerifyWebhook validates incoming webhook from payment provider
	VerifyWebhook(payload []byte, signature string) (bool, error)

	// GetPaymentStatus retrieves the current status of a payment
	GetPaymentStatus(transactionID string) (*PaymentResult, error)
}
