FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go mod tidy \
    && go build -o wowfishChain api/chain.go

FROM alpine:latest
WORKDIR /app/bin

# COPY --from=builder /app/api/etc/wowfishconfigTemp.yaml /app/bin/etc/wowfishconfig.yaml
COPY --from=builder /app/wowfishChain /app/bin/wowfishChain

# Copy start.sh script and define default command for the container
EXPOSE 8888
CMD [ "./wowfishChain" ]