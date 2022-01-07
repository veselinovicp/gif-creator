.DEFAULT_GOAL := help
.PHONY: help generate/fade gify/fade gify/mandelbrot clean generate/mandelbrot gif/mandelbrot
help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

generate/fade: ## Generate images
	go run fader/fader.go

gify/fade: ## Generate gif from images
	go run gifer/gifer.go -src="in_progress/*.png" -dest="in_progress/fade.gif" -scale=0.5 -verbose

gify/mandelbrot: ## Generate mandebrots from images
	go run gifer/gifer.go -src="in_progress/*.png" -dest="in_progress/mandelbrot22.gif" -scale=0.12 -delay=5 -verbose

generate/mandelbrot: ## Generate gif from images
	go run mandelbrot/mandelbrot.go -palette "Hippi" -xpos -0.00275 -ypos 0.78912 -radius .1256789 -iteration 800 -step 1000 -smoothness 15 -width 300 -height 300

gif/fade: generate/fade gify/fade ## Create faded gif

gif/mandelbrot: generate/mandelbrot gify/mandelbrot ## Create mandelbrot gif


#go run mandelbrot/mandelbrot.go -palette "Hippi" -xpos -0.00275 -ypos 0.78912 -radius .1256789 -iteration 800 -step 1000 -smoothness 15 -width 300 -height 300