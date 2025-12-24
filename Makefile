templ:
	@go tool templ generate -watch

air:
	@go tool air

tailwind:
	@npx tailwindcss -i input.css -o internal/adapter/web/static/css/styles.css --minify

dev:
	@make -j2 templ air

build:
	@go build -o bin/app cmd/server/main.go