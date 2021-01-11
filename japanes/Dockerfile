# build stage
FROM golang:alpine3.11 AS builder

ENV GO111MODULE=on
RUN apk --no-cache add tzdata
ADD ./ /go/src/japanes

RUN cd /go/src/japanes && go build -mod vendor

# final stage
FROM alpine:3.10.4

COPY --from=builder /go/src/japanes /go/src/japanes/.env.example /project/
### 下面這行是新加的設置時區
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai

WORKDIR /project

CMD ["./japanes", "-e", ".env.example"]