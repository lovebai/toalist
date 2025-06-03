FROM golang:alpine AS builder
WORKDIR /app
COPY . .
ARG GITHUB_SHA
RUN echo "Building commit: ${GITHUB_SHA:0:7}" && \
    go mod tidy && \
    go build -ldflags="-s -w -X main.CurrentCommit=${GITHUB_SHA:0:7}" -o main .

FROM alpine
ENV TZ=Asia/Shanghai
RUN apk add --no-cache alpine-conf ca-certificates  && \
    /usr/sbin/setup-timezone -z Asia/Shanghai && \
    apk del alpine-conf && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/main /app/main
CMD /app/main
EXPOSE 5245