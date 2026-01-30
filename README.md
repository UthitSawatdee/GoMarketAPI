# 🛒 E-Commerce REST API

A production-ready RESTful API for an e-commerce platform built with 
**Go (Golang)** following **Clean Architecture** principles.

## ✨ Features

- 🔐 **JWT Authentication** - Secure user authentication with role-based access control
- 👤 **User Management** - Registration, login, profile management, password change
- 📦 **Product Catalog** - CRUD operations with category management
- 🛒 **Shopping Cart** - Full cart functionality (add, update, remove, clear)
- 📋 **Order System** - Checkout, order tracking, cancellation
- 🏗️ **Clean Architecture** - Maintainable and testable codebase

## 🏗️ Architecture

This project follows **Clean Architecture** (Hexagonal Architecture) principles:

```
┌─────────────────────────────────────────┐
│           HTTP Handlers                 │  ← Adapters
├─────────────────────────────────────────┤
│           Use Cases                     │  ← Business Logic
├─────────────────────────────────────────┤
│           Entities                      │  ← Domain Models
├─────────────────────────────────────────┤
│        Repository (GORM)                │  ← Data Layer
└─────────────────────────────────────────┘
```
# รัน seed data
go run cmd/api/main.go -seed

### Public Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/register` | User registration |
| POST | `/api/v1/login` | User login |
| GET | `/api/v1/products` | List all products |
| GET | `/api/v1/products/:name` | Search products |
| GET | `/api/v1/products/category/:category` | Filter by category |

### Admin Endpoints (Admin Auth Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/admin/products` | Create product |
| PUT | `/api/v1/admin/products/:id` | Update product |
| DELETE | `/api/v1/admin/products/:id` | Delete product |
| POST | `/api/v1/admin/categories` | Create category |
| PUT | `/api/v1/admin/categories/:id` | Update category |
| DELETE | `/api/v1/admin/categories/:id` | Delete category |

### User Endpoints (Auth Required)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/user/profile` | Get profile |
| PUT | `/api/v1/user/profile` | Update profile/Change password  |