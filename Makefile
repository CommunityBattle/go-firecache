export GOOGLE_APPLICATION_CREDENTIALS=sa.json

.PHONY: test
test:
	go test
	
.PHONY: run
run:
	go run cmd/main.go