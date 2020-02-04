.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/upload upload/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/savenote savenote/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

validate: clean build
	sls deploy --noDeploy --stage dev

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
