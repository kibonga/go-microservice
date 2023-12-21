# name of the binaries (see dockerfile)
FRONTEND_BINARY=frontend
BROKER_BINARY=broker_app
AUTH_BINARY=auth_app
LOGGER_BINARY=logger_app

# starts all containers in the background without force build
up:
	@echo "Starting Docker images (detached, no-force)"
	docker-compose up
	@echo "Docker images started"

# stops docker-compose (if running), builds all projects and starts the containers in the background
up_build: build_broker build_auth build_logger
	@echo "Stopping docker images (if running)"
	docker-compose down
	@echo "Building (when required) and starting docker images (detached, force)"
	docker-compose up --build -d

# stops docker-compose
down:
	@echo "Stopping docker-compose"
	docker-compose down
	@echo "All containers stopped"

# builds the broker app binary as a linux executable
build_broker:
	@echo "Building broker app binary..."
	cd ../src/broker && env CGO_ENABLED=0 go build -o ${BROKER_BINARY} .
	@echo "Finished building broker app"

# builds the logger app binary as a linux executable
build_logger:
	@echo "Building logger app binary..."
	cd ../src/logger && env CGO_ENABLED=0 go build -o ${LOGGER_BINARY} .
	@echo "Finished building logger app"

build_auth:
	@echo "Building auth app binary..."
	cd ../src/auth && env CGO_ENABLED=0 go build -o ${AUTH_BINARY} .
	@echo "Finished building auth app"

# builds the frontend app as linux executable
build_front:
	@echo "Building frontend app binary..."
	cd ../src/frontend && env CGO_ENABLED=0 go build -o ${FRONTEND_BINARY} .
	@echo "Finished building frontend app"

# starts the frontend app
start_front:
	@echo "Starting frontend"
	cd ../src/frontend && "./${FRONTEND_BINARY}"

# stops the frontend app
stop_front:
	@echo "Stopping frontend"
	@-pkill -SIGTERM -f "./${FRONTEND_BINARY}"
	@echo "Stopped frontend"


