export GOPATH:=$(CURDIR)/_vendor:$(CURDIR)
export GOBIN:=$(CURDIR)/bin

all: install

fmt:
	gofmt -l -w -s src/

dep:
	go get github.com/go-sql-driver/mysql
	go get github.com/garyburd/redigo/redis
	go get github.com/wulijun/go-php-serialize/phpserialize
	go get github.com/gcliupeng/consumer/consumergroup
	go get github.com/Shopify/sarama

install: dep
	go install redirector

clean:
	rm _vendor/pkg -rf
	rm _vendor/src -rf
	rm bin/* -f
	rm pkg/* -rf
	rm log/* -f 

