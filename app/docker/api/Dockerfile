# build stage
FROM  golang:1.17-alpine3.16 AS Builder

WORKDIR /app

# ARG APP_CODE_PATH_HOST=../app

# RUN echo "I love my $APP_CODE_PATH_HOST.";


RUN pwd && ls -l
COPY ../../app/ .
# RUN go build -o main main.go


## run stage
# FROM alpine:3.16
# WORKDIR /app
## because builder main -o path is app/main
# COPY --from=Builder /app/main .
# COPY ${APP_CODE_PATH_HOST}/.env .

# CMD ["/app/main"]

