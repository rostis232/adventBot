BINARY_NAME=adventbot

build:
	@echo "Building..."
	go build -o ./bin/${BINARY_NAME} ./cmd/adventbot
	chmod +x ./bin/${BINARY_NAME}
	@echo "Built!"

run: build
	@echo "Starting..."
	./bin/${BINARY_NAME} &
	@echo "Started!"

stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./bin/${BINARY_NAME}"
	@echo "Stopped!"

start: run

restart: stop start