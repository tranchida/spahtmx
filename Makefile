MAKEFLAGS=-j 2

templ:
	@go tool templ generate -watch

air:
	@go tool air
    
dev:
	make -j2 templ air