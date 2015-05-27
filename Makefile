.PHONY: zip
zip:
	zip dynamic-dynamodb-lambda bootstrap.js dynamic-dynamodb-lambda

.PHONY: build
build: zip
	aws lambda create-function \
	--function-name dynamic-dynamodb-lambda \
	--zip-file "fileb://${PWD}/dynamic-dynamodb-lambda.zip" \
	--handler dynamic-dynamodb-lambda.handler \
	--runtime nodejs \
	--profile dynamic-dynamodb-lambda
