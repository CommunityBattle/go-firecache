export GOOGLE_APPLICATION_CREDENTIALS=../sa.json

.PHONY: test
test:
	go test ./test/e2e_test.go