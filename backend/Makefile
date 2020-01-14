run/dev:
	source .env && go run src/main.go

start/dev:
	source .env && \
	gin -a 8000 -p 8001 -i -t src/ run main.go

build/docker:
	docker build -t iamtheyammer/canvascbl/backend .
