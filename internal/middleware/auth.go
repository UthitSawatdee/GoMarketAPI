package middleware

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 1. ดึง Token จาก Header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "missing authorization header",
            })
        }

        // 2. ตรวจสอบ format "Bearer <token>"
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "invalid token format",
            })
        }

        // 3. Parse และ Validate Token
        token, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
            return []byte(secretKey), nil
        })

        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "invalid or expired token",
            })
        }

        // 4. ดึงข้อมูล user จาก token และเก็บใน context
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "invalid claims",
            })
        }

        // JWT decode ตัวเลขเป็น float64 ต้อง cast เป็น uint
        userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "user_id not found in token",
            })
        }

        // Set เป็น uint 
        c.Locals("userID", uint(userIDFloat))
        c.Locals("userRole", claims["role"])
        return c.Next()
    }
}
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง role จาก context (ที่ AuthMiddleware เก็บไว้)
		role := c.Locals("userRole")

		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "admin access required",
			})
		}

		return c.Next()
	}
}
