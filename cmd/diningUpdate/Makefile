SRC=build.go log.go main.go const.go
INTERNAL=../../internal/offers/*.go \
	../../internal/config/*.go \
	../../internal/timeout/*.go
EXE=diningUpdate
#GOLANG=/usr/local/go/bin/go
GOLANG=/usr/bin/go
CURL=/usr/bin/curl
GIT=/usr/bin/git
REPONAME=$(shell basename `git rev-parse --show-toplevel`)
SHA1=$(shell git rev-parse --short HEAD)
NOW=$(shell date +%Y-%m-%d_%T)
GOVERSION=$(shell go version | cut -f3 -d\ )
GOARCH=$(shell go version | cut -f4 -d\ )


diningUpdate: ${SRC} ${INTERNAL}
	echo ${REPONAME}
	${GOLANG} build \
      -ldflags "-X main.sha1ver=${SHA1} -X main.buildTime=${NOW} -X main.repoName=${REPONAME}"
	strip diningUpdate

all: update diningUpdate 

update: update-go

update-go:
	${GOLANG} get gopkg.in/ini.v1

pretty:
	${GOLANG} fmt ${SRC}

clean:
	rm -f ${EXE}

# vim: set noexpandtab:
