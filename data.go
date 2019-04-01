// > data.go

// WARNING DO NOT MANUALLY EDIT - YOUR CHANGES WILL BE OVERRIDDEN
// MAKE CHANGES AT ~/app/data/generate AND RUN make generate TO REGENERATE
// THE FOLLOWING FILE
//
// GENERATED BY GO:GENERATE AT 2019-04-01 23:13:50.175213187 +0800 HKT m=+0.029246505
//
// FILE GENERATED USING ~/app/data/generate.go

package main

// Version is used by godev for reporting the version when installed via 'go get'
const Version = "0.5.40"

// Commit is used by godev for reporting the version when installed via 'go get'
const Commit = "6444d53"

// DataDockerfile defines the 'Dockerfile' contents when --init is used
// hash:fc3c6491cb0d101ae17e2e68aec4714f
const DataDockerfile = `## 
## base image - defines the operating system layer for the build
## -------------------------------------------------------------
## use this to adjust the version of golang you want a build with
ARG GOLANG_VERSION=1.11.5
## use this to adjust the version of alpine to run for the build
ARG ALPINE_VERSION=3.9
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS base
## allow for passing in of any additional packages you might need
ARG ADDITIONAL_APKS
## due diligence
RUN apk update --no-cache
RUN apk upgrade --no-cache
## go modules dependencies
RUN apk add --no-cache git
## without these ssl/tls will not work
RUN apk add --no-cache ca-certificates && update-ca-certificates
## other development tooling
RUN apk add --no-cache make

##
## development image - where things are actually built
## ---------------------------------------------------
FROM base as development
## what should we name our binary? (default indicates "app")
ARG BIN_NAME=app
## any extension we would like for our binary? (default indicates nothing)
ARG BIN_EXT
## relative path to the binary from the working directory
ARG BIN_PATH=bin
## which architecture should we build for? (default indicates amd64)
ARG GOARCH=amd64
## which operating system should we build for? (default indicates linux)
ARG GOOS=linux
## should we use static linking? (default indicates yes)
ARG CGO_ENABLED=0
## should we use go modules for the dependencies? (default indicates yes)
ARG GO111MODULE=on
## use something GOPATH/GOROOT friendly - don't anger the gods
WORKDIR /go/src/${BIN_NAME}
## process dependencies first to take advantage of docker caching
COPY ./Makefile ./Makefile
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN make deps
## process everything else
COPY . /go/src/${BIN_NAME}
RUN make compile.linux
## generate a hash
RUN sha256sum ${BIN_PATH}/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} | cut -d " " -f 1 > ${BIN_PATH}/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256
## move things to where they should be
RUN mv /go/src/${BIN_NAME}/${BIN_PATH}/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} /${BIN_NAME}
RUN mv /go/src/${BIN_NAME}/${BIN_PATH}/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256 /${BIN_NAME}.sha256
RUN ln -s /${BIN_NAME} /_
RUN chmod +x /_
## let it start
ENTRYPOINT ["/_"]

##
# production image - the really small image
# -----------------------------------------
FROM scratch AS production
## what should we name our binary? (default indicates "app")
ARG BIN_NAME=app
WORKDIR /
## copy everything over from the development build image
COPY --from=base /etc/ssl/certs /etc/ssl/certs
COPY --from=development /${BIN_NAME} /${BIN_NAME}
COPY --from=development /${BIN_NAME}.sha256 /${BIN_NAME}.sha256
COPY --from=development /_ /_
## let it start
ENTRYPOINT ["/_"]
## if you're on openshift, you'll need to define this to define your application's ports
# EXPOSE 65534

`

// DataMakefile defines the 'Makefile' contents when --init is used
// hash:5ae20a94ab0d7cb71f695bdbf916bce9
const DataMakefile = `##
## Makefile constants - extract to a separate file if needed
## ---------------------------------------------------------
## specifies the name of your application binary
BIN_NAME=app
## specifies the relative path to a directory where the binary should be placed in
BIN_PATH=bin
## specifies the registry to push to
DOCKER_REGISTRY_HOSTNAME=docker.io
## specifies docker.io/THIS/image:tag
DOCKER_IMAGE_NAMESPACE=godev
## specifies docker.io/namespace/THIS:tag - align with $(BIN_NAME) for less confusion
DOCKER_IMAGE_NAME=example
## specifies the absolute path to the directory containing the .git directory
GIT_ROOT=$(CURDIR)
## enable following line to draw variables from a file named Makefile.properties
# include Makefile.properties

## starts the application for development with live-reload
start:
	@godev
## installs the dependencies using go modules
deps:
	@go mod vendor
## runs the tests with live-reload
test:
	@godev --test
## compiles binaries for all systems
compile:
	@$(MAKE) compile.linux
	@$(MAKE) compile.macos
	@$(MAKE) compile.windows
## compiles binaries for linux
compile.linux:
	@$(MAKE) GOARCH=amd64 GOOS=linux .compile
## compiles binaries for macos
compile.macos:
	@$(MAKE) GOARCH=amd64 GOOS=darwin .compile
## compiles binaries for windows
compile.windows:
	@$(MAKE) GOARCH=386 GOOS=windows BIN_EXT=.exe .compile
## compilation driver
.compile:
	@CGO_EMABLED=0 GO111MODULE=on \
		go build -a -ldflags "-extldflags -static" -o $(CURDIR)/$(BIN_PATH)/$(BIN_NAME)-${GOOS}-${GOARCH}${BIN_EXT}
	@chmod +x $(CURDIR)/$(BIN_PATH)/$(BIN_NAME)-${GOOS}-${GOARCH}${BIN_EXT}
	@sha256sum $(CURDIR)/$(BIN_PATH)/$(BIN_NAME)-${GOOS}-${GOARCH}${BIN_EXT} | cut -d " " -f 1 > $(CURDIR)/$(BIN_PATH)/$(BIN_NAME)-${GOOS}-${GOARCH}${BIN_EXT}.sha256
## dockerisation for production
docker:
	@$(MAKE) .docker STAGE="production"
## dockerisation for development
docker.dev:
	@$(MAKE) .docker STAGE="development"
## dockerisation driver
.docker:
	@$(MAKE) log.info MSG="creating image $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest"
	@docker build \
		--target ${STAGE} \
		--build-arg BIN_NAME=$(BIN_NAME) \
		--build-arg BIN_PATH=$(BIN_PATH) \
		--target=production \
		-t $(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		.
docker.prepare: docker
	@$(MAKE) log.info MSG="tagging image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest"
	@docker tag \
		$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest
	@$(MAKE) log.info MSG="tagging image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')"
	@docker tag \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')
	@$(MAKE) log.info MSG="tagging image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')-$$(git rev-list -1 HEAD)"
	@docker tag \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*') \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')-$$(git rev-list -1 HEAD)
publish.dockerhub: docker.prepare
	@$(MAKE) log.info MSG="pushing image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest"
	@docker push $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest
	@$(MAKE) log.info MSG="pushing image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')"
	@docker push $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')
	@$(MAKE) log.info MSG="pushing image $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')-$$(git rev-list -1 HEAD)"
	@docker push $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):$$($(MAKE) version.get | grep '[0-9]*\.[0-9]*\.[0-9]*')-$$(git rev-list -1 HEAD)
version.get:
	@docker run \
		-v "$(GIT_ROOT):/app" \
		zephinzer/vtscripts:latest \
		get-latest -q
version.next:
	@docker run \
		-v "$(GIT_ROOT):/app" \
		zephinzer/vtscripts:latest \
		get-next -q
version.bump:
	@docker run \
		-v "$(GIT_ROOT):/app" \
		zephinzer/vtscripts:latest \
		iterate ${VERSION} -i -q
log.debug:
	-@printf -- "\033[36m\033[1m_ [DEBUG] ${MSG}\033[0m\n"
log.info:
	-@printf -- "\033[32m\033[1m>  [INFO] ${MSG}\033[0m\n"
log.warn:
	-@printf -- "\033[33m\033[1m?  [WARN] ${MSG}\033[0m\n"
log.error:
	-@printf -- "\033[31m\033[1m! [ERROR] ${MSG}\033[0m\n"

`

// DataDotGitignore defines the '.gitignore' contents when --init is used
// hash:3e59a1165602d77a63163af48e9793bc
const DataDotGitignore = `# development artifacts
bin
c.out
vendor

`

// DataDotDockerignore defines the '.dockerignore' contents when --init is used
// hash:9441e48bcf7b0249fc852973e74053f4
const DataDotDockerignore = `.dockerignore
.gitignore
Dockerfile
bin
c.out
vendor

`

// DataMainDotgo defines the '.dockerignore' contents when --init is used
// hash:4a73f12d9bde8b278abb6dc558584402
const DataMainDotgo = `package main

import "fmt"

func main() {
	fmt.Println("hello world!")
}

`

// DataGoDotMod defines the 'go.mod' contents when --init is used
// hash:b6791696ce7f0e334775b206f1fa9dea
const DataGoDotMod = `module app

`


// WARNING DO NOT MANUALLY EDIT - YOUR CHANGES WILL BE OVERRIDDEN
// MAKE CHANGES AT ~/app/data/generate AND RUN make generate TO REGENERATE
// THE FOLLOWING FILE
//
// GENERATED BY GO:GENERATE AT 2019-04-01 23:13:50.175213187 +0800 HKT m=+0.029246505
//
// FILE GENERATED USING ~/app/data/generate.go

// < data.go
