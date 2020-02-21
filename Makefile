build:
	@go build

test:
	@go test -race -p 1 -failfast -v ./...

lint:
	@golangci-lint run --enable-all
	@find . -iname "*.sh" -exec shellcheck -x {} +
	@find . \( -iname "*.yml" -o -iname "*.yaml" -o -iname "*.md" -o -iname "*.json" -o -iname ".prettierrc" \) -exec npx prettier -c {} +

format:
	@golangci-lint run --enable-all --fix
	@find . \( -iname "*.yml" -o -iname "*.yaml" -o -iname "*.md" -o -iname "*.json" -o -iname ".prettierrc" \) -exec npx prettier --write {} +