PLATFORM=$(shell uname -m)
DATETIME=$(shell date "+%Y%m%d%H%M%S")
VERSION=v0.0.1-SNAPSHOT

#default: backend-debugging

# release:
#	@cd front && npm run build -- prod && cp -af dist/* ../service/src/resources/
#	@cd src && GOPATH=${GOPATH} go build -o ../bin/server.bin

#front-debugging:
#	@cd front && npm run dev

backend-debugging:
	@cd src && GOPATH=${GOPATH} go build -gcflags=all="-N -l" -o ../bin/search.bin

clean:
	@rm -rf ./bin

#docker:
#	@docker build . -t harbor.timeforward.cn:8443/public/bleve-searcher:$(VERSION)-$(DATETIME)

