SRC=build.go log.go main.go const.go
INTERNAL=../../internal/offers/offers.go ../../internal/mail/mail.go
EXE=diningUpdate
#GOLANG=/usr/local/go/bin/go
GOLANG=/usr/bin/go
CURL=/usr/bin/curl
GIT=/usr/bin/git
REPONAME=$(shell basename `git rev-parse --show-toplevel`)
SHA1=$(shell git rev-parse --short HEAD)
NOW=$(shell date +%Y-%m-%d_%T)
#export HTTP_PROXY=http://llproxy.llan.ll.mit.edu:8080
#export HTTPS_PROXY=${HTTP_PROXY}

diningUpdate: ${SRC} ${INTERNAL}
	echo ${REPONAME}
	${GOLANG} build \
      -ldflags "-X main.sha1ver=${SHA1} -X main.buildTime=${NOW} -X main.repoName=${REPONAME}"

update: update-go update-json

update-go:
	${GOLANG} get gopkg.in/ini.v1

pretty:
	${GOLANG} fmt ${SRC}

clean:
	rm -f ${EXE}

# vim: set noexpandtab:
