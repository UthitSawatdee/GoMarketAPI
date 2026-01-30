package handler

import (
	// domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	usecases "github.com/Fal2o/E-Commerce_API/internal/usecases"
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
// @Description Cancel a specific order by its ID
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param orderID path string true "Order ID"
// @Success 200 {object} map[string]interface{} "Order canceled successfully"
// @Failure 400 {object} map[string]interface{} "Order ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/order/cancel/{orderID} [delete]
func (h *HttpOrderHandler) CancelOrder(c *fiber.Ctx) error {
	// Implementation for canceling an order goes here
	orderID := c.Params("orderID")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	err := h.OrderUseCase.CancelOrder(orderID)
	if err != nil {

		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel order",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order canceled successfully",
	})
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
// @Description Update the status of a specific order (Admin only)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param orderID path string true "Order ID"
// @Param status path string true "New status" Enums(pending, processing, shipped, delivered, cancelled)
// @Success 200 {object} map[string]interface{} "Order status updated successfully"
// @Failure 400 {object} map[string]interface{} "Order ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/order/status/{orderID}/{status} [put]
func (h *HttpOrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	// Implementation for updating order status goes here
	orderID := c.Params("orderID")
	status := c.Params("status")

	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	order, oldStatus, err := h.OrderUseCase.UpdateOrderStatus(orderID, status)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
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
