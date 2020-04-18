FROM golang:1.12-buster as builder

ENV GOOS=linux \
    GOARCH=amd64

RUN mkdir -p /opt/app && \
    cd /opt/app && \
    apt update && \
    apt install -y --no-install-recommends xz-utils && \
    go get -d -v gopkg.in/gomail.v2 github.com/gin-gonic/gin && \
    go get -v github.com/pwaller/goupx && \
    curl -ksSL https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz -o upx.tar.xz && \
    mkdir upx && \
    tar -xf upx.tar.xz --strip-components 1 -C upx

WORKDIR /opt/app

COPY *.go /opt/app/

RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o mailer.bin . && \
    PATH="$PWD/upx:$PATH" goupx mailer.bin

FROM scratch

COPY --from=builder /opt/app/mailer.bin /bin/mailer

CMD [ "mailer", "--http" ]
