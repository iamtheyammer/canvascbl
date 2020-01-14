.PHONY: frontend/build frontend/formatcheck backend/build

frontend/build:
	if [ ! -d "frontend/node_modules" ]; then cd frontend && npm install; fi;
	cd frontend && npm run build;
	mv ./frontend/build ./bin;

frontend/formatcheck:
	cd frontend && npm run formatcheck;

backend/build:
	cd backend && go build -o ../bin/canvasProxy src/main.go

ci:
	make backend/build;
	make frontend/formatcheck;
	make frontend/build;
