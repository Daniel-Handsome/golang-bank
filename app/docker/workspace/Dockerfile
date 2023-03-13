ARG GOLANG_VERSION
FROM  golang:${GOLANG_VERSION}-alpine3.16 
RUN go get -u github.com/cosmtrek/air
RUN apk add alpine-sdk
WORKDIR /app
ENTRYPOINT ["air"]