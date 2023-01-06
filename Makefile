#aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name sandman-queue
#
##sudo snap install aws-cli --classic
#create-local-queue:
# 	@aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name sandman-queue

build:
	@go build -o app cmd/main.go #gosetup


.PHONY: test
test:
	@go test -json ./... -covermode=atomic -coverprofile coverage.out
#	@go test -c -coverpkg=... -covermode=atomic -o test.test lab.dev.vm.co.mz/compse/sandman
#	@go tool test2json -t ./test.test -test.v -test.coverprofile coverage.out

upload-pkg:


download-pkg:
	@curl --user 9YM0pOGm:HKEohJeeC4dD4GIKcvetJfoyh3kbt_LzO229JltQSIOY -o 1.0.0-SNAPSHOT.25.zip https://nexus.pkg.dev.vm.co.mz/repository/zip-public/commission-engine-payment-lambda/1.0.0-SNAPSHOT.25.zip
