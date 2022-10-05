FROM golang:1.19-alpine
RUN go install github.com/beego/bee/v2@latest
RUN apk add build-base
ENV GO111MODULE=on
ENV APP_HOME /go/src/barckend
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"
COPY ./ ./
EXPOSE 8000
CMD ["bee", "run"]
