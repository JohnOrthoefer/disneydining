SRC= main.go build.go startWatcher.go
INTERNAL=../../internal/config/*.go
EXE=diningWeb
GOLANG=/usr/local/go/bin/go
#GOLANG=/usr/bin/go
CURL=/usr/bin/curl
GIT=/usr/bin/git
REPONAME=$(shell basename `git rev-parse --show-toplevel`)
SHA1=$(shell git rev-parse --short HEAD)
NOW=$(shell date +%Y-%m-%d_%T)
GOVERSION=$(shell go version | cut -f3 -d\ )
GOARCH=$(shell go version | cut -f4 -d\ )

diningWeb: ${SRC} ${INTERNAL}
	echo ${REPONAME}
	${GOLANG} build \
      -ldflags "-X main.sha1ver=${SHA1} -X main.buildTime=${NOW} -X main.repoName=${REPONAME} -X main.goVersion=${GOVERSION} -X main.goArch=${GOARCH}"
	strip diningWeb

all: update diningWeb

update: update-go

update-go:
	${GOLANG} get gopkg.in/ini.v1


pretty:
	${GOLANG} fmt ${src}

clean:
	rm -f ${EXE}

# vim: set noexpandtab:
