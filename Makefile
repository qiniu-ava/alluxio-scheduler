SHELL:=/bin/bash
GOROOT:=$(shell echo ${GOROOT})
GOPATH:=$(shell pwd)/app
DIR:=$(shell pwd)

clean:
	rm ${DIR}/app/bin/server

build_debug:
	cd ${DIR}/app/src/qiniu.com/server && go install -tags debug

build_prd:
	cd ${DIR}/app/src/qiniu.com/server && go install

build: build_prd

deploy:
	${DIR}/docker-build.sh
