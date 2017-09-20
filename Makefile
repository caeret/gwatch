BINARY_FILE=gwatch
GO=$(shell which go)

build: getdeps
	@echo "Running $@" && ${GO} build -o ${BINARY_FILE} main.go

install: getdeps
	@echo "Running $@" && ${GO} install

getdeps:
	@echo "Installing dependencies by glide" && glide -debug install

clean:
	@echo "Cleaning up all the generated files"
	@rm -f ${BINARY_FILE}
	@rm -rf vendor

.PHONY: get clean

