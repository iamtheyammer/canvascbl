frontendbuild:
	if [ ! -d "frontend/node_modules" ]; then cd frontend && npm install; fi;
	cd frontend && npm run build;
	mv ./frontend/build ./bin;

build:
	go build -o bin/canvasProxy src/main.go

devrunbuilt: build;
	source .env && ./bin/canvasProxy

devrun:
	source .env && go run src/main.go

devstart:
	cd src && \
	source ../.env && \
	gin -a 8000 -p 8001 -i run main.go

ci:
	make build;
	cd frontend && npm run formatcheck && cd ..;
	make frontendbuild;
