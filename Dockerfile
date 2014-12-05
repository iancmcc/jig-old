FROM       ubuntu:trusty
MAINTAINER Ian McCracken <ian.mccracken@gmail.com>

RUN apt-get update && apt-get install -y build-essential curl

# Install Go
RUN   curl -sSL https://golang.org/dl/go1.3.3.src.tar.gz | tar -v -C /usr/local -xz
ENV   PATH    /usr/local/go/bin:$PATH
ENV   GOPATH  /go:/go/src/github.com/iancmcc/jig/Godeps/_workspace
ENV   PATH /go/bin:$PATH
RUN   cd /usr/local/go/src && ./make.bash --no-clean 2>&1

# Compile Go for cross compilation
 ENV   DOCKER_CROSSPLATFORMS   \
       linux/386 linux/arm \
       darwin/amd64 darwin/386 \
       freebsd/amd64 freebsd/386 freebsd/arm 

ENV GOARM 5
RUN cd /usr/local/go/src && bash -xc 'for platform in $DOCKER_CROSSPLATFORMS; do GOOS=${platform%/*} GOARCH=${platform##*/} ./make.bash --no-clean 2>&1; done'

WORKDIR /go/src/github.com/iancmcc/jig
COPY  .   /go/src/github.com/docker/docker