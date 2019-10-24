herokunode:
	if [ ! -d "frontend/node_modules" ]; then cd frontend && yarn; fi;
	cd frontend && yarn build;
	mv ./frontend/build ./bin;

build:
	go build -o bin/canvasProxy src/main.go

devrunbuilt: build;
	source .env && ./bin/canvasProxy

devrun:
	source .env && go run src/main.go

ci:
	make build;
	cd frontend && yarn formatcheck && cd ..;
	make herokunode;
