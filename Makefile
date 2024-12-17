.DEFAULT_GOAL := help

##### BUILD #####

.PHONY: build/server
build/server: ## Build the server
	go build -o .artifacts/server ./app

.PHONY: run/queue-sample
run/queue-sample: ## Run the queue sample
	go run app/main.go queue-sample

.PHONY: run/gen-migration
run/gen-migration: ## generate new migration up/down files ## make run/gen-migration n=create_user_table
run/gen-migration:
	go run ./app/main.go gen-migration -n $(n)

.PHONY: help
help: ## Display this help screen ## make or make help
	@echo ""
	@echo "Usage: make SUB_COMMAND argument_name=argument_value"
	@echo ""
	@echo "Command list:"
	@echo ""
	@printf "\033[36m%-30s\033[0m %-50s %s\n" "[Sub command]" "[Description]" "[Example]"
	@grep -E '^[/a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | perl -pe 's%^([/a-zA-Z_-]+):.*?(##)%$$1 $$2%' | awk -F " *?## *?" '{printf "\033[36m%-30s\033[0m %-50s %s\n", $$1, $$2, $$3}'