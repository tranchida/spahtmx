templ:
	@go tool templ generate -watch

air:
	@go tool air
    
dev:
	@make -j2 templ air

build:
	@go build -o bin/app main.go