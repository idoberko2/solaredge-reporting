BINARY_NAME=semonitor
BINARY_DIR=bin

build:
	go build -o ${BINARY_DIR}/${BINARY_NAME} main.go
	GOARCH=arm64 GOOS=linux go build -o ${BINARY_DIR}/${BINARY_NAME}-linux main.go

test:
	go test ./...

test-integration:
	go test ./... -tags=integration

deploy: build
	scp -i ~/.ssh/personal_rsa bin/semonitor-linux idoberko2@raspberrypi.local:~/
	scp -i ~/.ssh/personal_rsa db/migrations/* idoberko2@raspberrypi.local:~/db/migrations

run: build
	./${BINARY_DIR}/${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_DIR}/${BINARY_NAME}
	rm ${BINARY_DIR}/${BINARY_NAME}-linux