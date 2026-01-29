.PHONY: seed

seed:
	@echo " Seeding database..."
	@go run cmd/api/main.go -seed
	@echo " Seed complete"