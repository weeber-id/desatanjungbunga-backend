FROM golang:1.14.4-alpine3.12 as build-go-deps
WORKDIR /work
COPY ./ ./
RUN go build -o main.out main.go


FROM alpine:3.12.0

RUN apk add tzdata
ENV TZ=Asia/Jakarta
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV GIN_MODE=release

WORKDIR /app
COPY --from=build-go-deps /work/main.out /app
COPY --from=build-go-deps /work/VERSION /app
EXPOSE 8080
CMD ["./main.out"]