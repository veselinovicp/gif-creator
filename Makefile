.DEFAULT_GOAL := help
.PHONY: help generate gify clean
help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

generate: ## Generate images
	$(MAKE) -C ./code/go generate

gify: ## Generate gif from images
	$(MAKE) -C ./code/go gify

gif/fade:
	$(MAKE) -C ./code/go gif/fade

clean: ## Clean intermediate files
	rm -rf code/go/in_progress