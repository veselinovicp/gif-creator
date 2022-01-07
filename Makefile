.DEFAULT_GOAL := help
.PHONY: help generate/fade gify clean generate/mandelbrot gif/mandelbrot
help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

generate/fade: ## Generate images
	$(MAKE) -C ./code/go generate/fade

gify: ## Generate gif from images
	$(MAKE) -C ./code/go gify

generate/mandelbrot: 
	$(MAKE) -C ./code/go generate/mandelbrot

gif/fade:
	$(MAKE) -C ./code/go gif/fade

gif/mandelbrot:
	$(MAKE) -C ./code/go gif/mandelbrot

clean: ## Clean intermediate files
	rm -rf code/go/in_progress