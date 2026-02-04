package domain

// OrderStatus represents the status of an order in its lifecycle
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// ValidTransitions defines allowed state transitions for orders
// This enforces business rules about order lifecycle
var ValidTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending:    {OrderStatusConfirmed, OrderStatusCancelled},
	OrderStatusConfirmed:  {OrderStatusProcessing, OrderStatusCancelled},
	OrderStatusProcessing: {OrderStatusShipped, OrderStatusCancelled},
	OrderStatusShipped:    {OrderStatusDelivered},
	OrderStatusDelivered:  {}, // Terminal state - no further transitions allowed
	OrderStatusCancelled:  {}, // Terminal state - no further transitions allowed
}

// CanTransitionTo checks if the current status can transition to the target status
func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	allowed, exists := ValidTransitions[s]
	if !exists {
		return false
	}
	for _, v := range allowed {
		if v == next {
			return true
		}
	}
	return false
}

// IsValid checks if the status is a recognized order status
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusProcessing,
		OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled:
		return true
	}
	return false
}

// IsCancellable checks if an order with this status can be cancelled
func (s OrderStatus) IsCancellable() bool {
	return s.CanTransitionTo(OrderStatusCancelled)
}

// IsTerminal checks if the status is a terminal state (no further changes allowed)
func (s OrderStatus) IsTerminal() bool {
	return s == OrderStatusDelivered || s == OrderStatusCancelled
}
