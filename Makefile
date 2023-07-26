reset-storage:
	@echo "Resetting storage..."
	@rm -rf ./storage
	@mkdir ./storage
	@mkdir ./storage/0
	@mkdir ./storage/1
	@mkdir ./storage/2
	@mkdir ./storage/3
	@mkdir ./storage/4
	@echo "Done."

test:
	@echo "Running tests..."
	@./smoke_test.sh
	@echo "Done."

run:
	env $$(cat .env) go run cmd/main.go