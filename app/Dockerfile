## 放這裡的原因是 copy只能往下 不能../這樣 所以他一定要再上層
################  build stage
FROM  golang:1.17-alpine3.16 AS Builder

LABEL maintainer="daniel <adwxsghu@gmail.com>"

WORKDIR /app

COPY . .
RUN go build -o main main.go

# RUN apk update && apk upgrade && \
#     apk --update --no-cache add curl && \
#     apk --update --no-cache add tar && \
#     apk --update --no-cache add gzip

# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.0/migrate.linux-amd64.tar.gz | tar xvz



################ run stage
FROM alpine:3.16
WORKDIR /app
# because builder main -o path is app/main
COPY --from=Builder /app/main .
COPY .env .
COPY start.sh .

### migrate set
# COPY migrate.linux-amd64  ./migrate
# COPY --from=Builder /app/migrate.linux-amd64 ./migrate
## 複製外部的進來給他跑migrate 路徑最好一樣
COPY db/migrations ./db/migrations

## RUN sh 

## ENTRYPOINT cme會當變數 ENTRYPOINT
## 當然可以放在entrypoint的第二的餐數

ARG POSTGRESQL_URL

EXPOSE 8080

ENTRYPOINT [ "/app/main" ]
# CMD /app/main
# ENTRYPOINT [ "chmod", "+x", "/app/start.sh" ]
# CMD ["/app/main"]

