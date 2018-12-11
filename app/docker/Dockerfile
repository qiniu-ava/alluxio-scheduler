FROM golang:1.10.5-alpine as builder
WORKDIR /go/src
ARG VERSION=prd
COPY ./ /go/src/
RUN cd qiniu.com/server && go install

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/server /
COPY --from=builder /go/src/qiniu.com/server/config/config.prd.json /config/config.prd.json
COPY --from=builder /go/src/qiniu.com/server/config/config.default.json /config/config.default.json
CMD ["/server"]
