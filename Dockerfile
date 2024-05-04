FROM golang:1.22.2-alpine3.19

WORKDIR /app/bin

RUN apk add --no-cache curl tar && \
   curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip && \
   unzip protoc-3.17.3-linux-x86_64.zip -d /usr/local && \
   rm -f protoc-3.17.3-linux-x86_64.zip

COPY ./ ./

RUN go mod tidy

RUN go build main.go

CMD ["./main"]