package handler

import (
	// domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	usecases "github.com/UthitSawatdee/GoMarketAPI/internal/usecases"
	"github.com/gofiber/fiber/v2"
)

type HttpOrderHandler struct {
	OrderUseCase usecases.OrderUseCase
}

func NewHttpOrderHandler(useCase usecases.OrderUseCase) *HttpOrderHandler {
	return &HttpOrderHandler{OrderUseCase: useCase}
}

// ViewOrder godoc
// @Summary View user orders
// @Description Get all orders for the authenticated user
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Orders retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/orders [get]
func (h *HttpOrderHandler) ViewOrder(c *fiber.Ctx) error {
	// Implementation for viewing orders goes here
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	request, err := h.OrderUseCase.ViewOrder(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Orders retrieved successfully",
		"data":    request,
	})
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Cancel a specific order by its ID. Stock will be restored to inventory.
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Order canceled successfully"
// @Failure 400 {object} map[string]interface{} "Order ID is required or order cannot be cancelled"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Cannot cancel another user's order"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/order/cancel/{orderID} [delete]
func (h *HttpOrderHandler) CancelOrder(c *fiber.Ctx) error {
	orderID := c.Params("orderID")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	// Get userID from JWT context for ownership verification
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err := h.OrderUseCase.CancelOrder(orderID, userID)
	if err != nil {
		errMsg := err.Error()

		// Handle specific error cases
		if contains(errMsg, "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if contains(errMsg, "unauthorized") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Cannot cancel another user's order",
			})
		}
		if contains(errMsg, "cannot be cancelled") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel order",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order canceled successfully. Stock has been restored.",
	})
}

// contains checks if substr is in s (helper function)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ViewAllOrders godoc
// @Summary Get all orders
// @Description Retrieve all orders in the system (Admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "All orders retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/orders [get]
func (h *HttpOrderHandler) ViewAllOrders(c *fiber.Ctx) error {
	// Implementation for viewing all orders goes here
	request, err := h.OrderUseCase.AllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve orders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All orders retrieved successfully",
		"data":    request,
	})
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of a specific order (Admin only). Validates state transitions.
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param orderID path string true "Order ID"
// @Param status path string true "New status" Enums(pending, confirmed, processing, shipped, delivered, cancelled)
// @Success 200 {object} map[string]interface{} "Order status updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid status or transition"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/order/status/{orderID}/{status} [put]
func (h *HttpOrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("orderID")
	status := c.Params("status")

	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Status is required",
		})
	}

	order, oldStatus, err := h.OrderUseCase.UpdateOrderStatus(orderID, status)
	if err != nil {
		errMsg := err.Error()

		if contains(errMsg, "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if contains(errMsg, "invalid transition") || contains(errMsg, "invalid status") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Order status updated successfully",
		"old_status": oldStatus,
		"new_status": order.Status,
	})
}
