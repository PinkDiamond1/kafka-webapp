FROM fedora:latest

LABEL maintainer="Amit Kumar Gupta <amitkgupta84@gmail.com>"

RUN dnf install wget -y
RUN dnf install gcc -y

RUN rpm --import https://packages.confluent.io/rpm/5.4/archive.key
COPY confluent.repo /etc/yum.repos.d
RUN dnf install librdkafka-devel -y

RUN wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz && tar -xvf go1.14.linux-amd64.tar.gz && rm go1.14.linux-amd64.tar.gz
RUN mv go /usr/local
ENV GOROOT=/usr/local/go
ENV PATH="${GOROOT}/bin:${PATH}"

WORKDIR /
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY produce.go produce.go
COPY consume.go consume.go
RUN go build -o server
RUN rm *.go && rm go.*

COPY assets assets
COPY produce.html produce.html
COPY consume.html consume.html

ENV PORT 8080
EXPOSE 8080
ENTRYPOINT ["/server"]
