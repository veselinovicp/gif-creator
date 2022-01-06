.DEFAULT_GOAL := help
.PHONY: help generate gify clean
help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

generate: ## Generate images
	go run fader/fader.go

gify: ## Generate gif from images
	go run gifer/gifer.go -src="in_progress/*.png" -dest="in_progress/movie.gif" -scale=0.5 -verbose

gif/fade: generate gify ## Create faded gif
