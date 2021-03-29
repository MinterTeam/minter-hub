all: install

install: go.sum
		go install $(BUILD_FLAGS) ./cmd/mhub-keys-generator

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify
