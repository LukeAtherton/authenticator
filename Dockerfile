FROM golang
MAINTAINER Luke Atherton "lukeatherton@hivebase.io"

# Install dependencies
run apt-get update
run apt-get install -y curl git

# Install confd
RUN curl -L https://github.com/kelseyhightower/confd/releases/download/v0.3.0/confd_0.3.0_linux_amd64.tar.gz | tar -xz
RUN mv confd /usr/local/bin/confd

RUN go get github.com/tools/godep

RUN mkdir -p /etc/confd/conf.d
RUN mkdir -p /etc/confd/templates

WORKDIR /go/src/github.com/lukeatherton/authenticator
ADD . /go/src/github.com/lukeatherton/authenticator

# Add files
ADD ./bin/boot.sh             /boot.sh
ADD ./confd                   /etc/confd
ADD ./crypto                  /crypto

EXPOSE 8001

# Start the container
RUN chmod +x /boot.sh
CMD /boot.sh

#Install App
RUN godep go build
