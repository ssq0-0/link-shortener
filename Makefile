BINARY_NAME=Linker

.PHONY:
all: build

.PHONY:
build: 
	go build -o ${BINARY_NAME}

.PHONY:
run: build
	./${BINARY_NAME}

