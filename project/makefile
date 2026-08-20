SHELL=cmd.exe
FRONT_END_BINARY=frontApp.exe

# up: starts all containers in the background without forcing build
up:
	@echo Starting all Docker images
	docker-compose up -d
	@echo All Docker images started!

# up_build: stops docker-compose (if running), rebuilds all Docker \
images and starts docker compose
up_build:
	@echo Stoping all Docker images (if running)...
	docker-compose down
	@echo Rebuilding all Docker images and starting the built Docker images
	docker-compose --build -d
	@echo All Docker images started!

# down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

# build_front: builds the front end binary
build_front:
	@echo Building the front end binary
	cd ..\front-end && SET CGO_ENABLED=0 && go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo Done!

#start: starts the front end
start:
	@echo Starting the front end
	cd ./../front-end && start /B ${FRONT_END_BINARY}
	@echo Front end started

#stop: stop the front end
stop:
	@echo Stopping the front end...
	@taskill /IM "${FRONT_END_BINARY}" /F
	@echo The front end is stopped!

