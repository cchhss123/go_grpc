FROM golang:1.23-alpine AS builder

RUN apk update && apk add --no-cache git build-base protobuf wget

RUN wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v30.2/protoc-30.2-linux-x86_64.zip
RUN unzip protoc.zip -d /usr/local
RUN rm protoc.zip

# 設置 GOBIN
ENV GOBIN=/go/bin

# 安装 protoc-gen-go 和 protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# GOBIN 和 protoc 的安装目录 加到 PATH 中
ENV PATH="$PATH:$GOBIN:/usr/local/bin"


# 複製當前目錄的程式碼到容器內的工作目錄
COPY ./app /app

# 設定工作目錄
WORKDIR /app


