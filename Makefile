test: 
	@go test -failfast -count=1 ./... -json -cover -race | tparse -smallscreen