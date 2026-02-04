package migrations

import (
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	"gorm.io/gorm"
	"log"
	hash "github.com/UthitSawatdee/GoMarketAPI/pkg/hash"
)

// SeedData inserts initial data into database
func SeedData(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	// 1. Seed Categories
	categories := []domain.Category{
		{Name: "Electronics", Description: "Electronic devices and gadgets"},
		{Name: "Fashion", Description: "Clothing and accessories"},
		{Name: "Home & Garden", Description: "Home decor and garden tools"},
		{Name: "Books", Description: "Physical and digital books"},
	}

	for _, category := range categories {
		// Check if category already exists
		var existing domain.Category
		if err := db.Where("name = ?", category.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			// Category doesn't exist, create it
			if err := db.Create(&category).Error; err != nil {
				log.Printf("Failed to seed category %s: %v", category.Name, err)
			} else {
				log.Printf("Seeded category: %s", category.Name)
			}
		}
	}

	// 2. Get category IDs for products
	var electronics, fashion domain.Category
	db.Where("name = ?", "Electronics").First(&electronics)
	db.Where("name = ?", "Fashion").First(&fashion)

	// 3. Seed Products
	products := []domain.Product{
		{
			Name:        "iPhone 15 Pro",
			Description: "Latest Apple smartphone with A17 Pro chip",
			Price:       999.99,
			Stock:       50,
			CategoryID:  electronics.ID,
			// IsActive:    true,
		},
		{
			Name:        "Samsung Galaxy S24",
			Description: "Flagship Android phone with AI features",
			Price:       899.99,
			Stock:       30,
			CategoryID:  electronics.ID,
			// IsActive:    true,

		},
		{
			Name:        "Nike Air Max",
			Description: "Comfortable running shoes",
			Price:       129.99,
			Stock:       100,
			CategoryID:  fashion.ID,
			// IsActive:    true,

		},
	}

	for _, product := range products {
		// Check if product already exists
		var existing domain.Product
		if err := db.Where("name = ?", product.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&product).Error; err != nil {
				log.Printf("Failed to seed product %s: %v", product.Name, err)
			} else {
				log.Printf("Seeded product: %s", product.Name)
			}
		}
	}

	users := []domain.User{
		{
			Username: "AdminName",
			Email:    "admin@admin.com",
			Password: "mypassword",
			Role:     "admin",
		},
		{
			Username: "UserName",
			Email:    "user@user.com",
			Password: "mypassword",
			Role:     "user",
		},
	}

    for _, user := range users {
        var existing domain.User
        if err := db.Where("email = ?", user.Email).First(&existing).Error; err == gorm.ErrRecordNotFound {
            // ✅ Hash password ก่อน save
            hashedPassword, err := hash.NewPasswordService().Hash(user.Password)
            if err != nil {
                log.Printf(" Failed to hash password for %s: %v", user.Email, err)
                continue
            }
            
            user.Password = hashedPassword // Set hashed password
            
            if err := db.Create(&user).Error; err != nil {
                log.Printf(" Failed to seed user %s: %v", user.Email, err)
            } else {
                log.Printf(" Seeded user: %s (Role: %s)", user.Email, user.Role)
            }
        } else {
            log.Printf(" User %s already exists, skipping", user.Email)
        }
    }

	log.Println("Seeding completed!")
	return nil
}
