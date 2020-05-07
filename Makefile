build:
	@go build

test:
	@AWS_ACCESS_KEY_ID=makefile AWS_SECRET_ACCESS_KEY=localstack go test -cover -race -p 1 -failfast -v ./...

lint:
	@golangci-lint run --enable-all
	@find . -iname "*.sh" -exec shellcheck -x {} +
	@find . \( -iname "*.yml" -o -iname "*.yaml" -o -iname "*.md" -o -iname "*.json" -o -iname ".prettierrc" \) -exec npx prettier -c {} +

format:
	@golangci-lint run --enable-all --fix
	@find . \( -iname "*.yml" -o -iname "*.yaml" -o -iname "*.md" -o -iname "*.json" -o -iname ".prettierrc" \) -exec npx prettier --write {} +

security:
	@yarn audit
	@snyk test --all-projects --dev
