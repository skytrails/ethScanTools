FROM golang:alpine as builder

MAINTAINER flynn

ENV GOPROXY https://goproxy.cn/

WORKDIR /go/release
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add tzdata

COPY go.mod ./go.mod
RUN go mod tidy
COPY . .
RUN pwd && ls

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o eth-scan .

FROM alpine

COPY --from=builder /go/release/eth-scan /

#COPY --from=builder /go/release/config/settings.gen.yml /config/settings.gen.yml
#COPY --from=builder /go/release/config/settings.scan.yml /config/settings.scan.yml

#COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

EXPOSE 8000

#CMD ["/eth-scan","scan","-c", "/config/settings.scan.yml"]
CMD ["/eth-scan","gen","-c", "/config/settings.gen.yml"]
