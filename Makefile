default: install-content-checker

.PHONY: install-content-checker
install-content-checker:
	go build -o ${GOPATH}/bin/mdcheck ./src/a1.2/go-essentials/2-content-checker
