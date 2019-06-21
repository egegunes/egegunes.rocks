run:
	go run main.go
rundb:
	docker run -d --rm --name dynamodb -p 8000:8000 amazon/dynamodb-local
