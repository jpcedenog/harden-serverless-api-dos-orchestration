.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/guesssecret guesssecret/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

remove:
	sls remove

tfinit:
	terraform init

tfapply:
	terraform apply

tfdestroy:
	terraform destroy
