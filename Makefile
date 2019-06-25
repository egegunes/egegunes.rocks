all: build package deploy clean

build:
	go build -o egegunesrocks
package:
	zip -r egegunesrocks.zip egegunesrocks .env templates/
deploy:
	aws cloudformation package --template-file sam.yaml --output-template-file output-sam.yaml --s3-bucket egegunesrocks
	aws cloudformation deploy --template-file output-sam.yaml --stack-name egegunesrocks --capabilities CAPABILITY_IAM
clean:
	rm egegunesrocks egegunesrocks.zip output-sam.yaml
run:
	go run main.go
rundb:
	docker run -d --rm --name dynamodb -p 8000:8000 amazon/dynamodb-local
