ARG GOLANG_VERSION
FROM  golang:${GOLANG_VERSION}-alpine3.16 
RUN go get -u github.com/cosmtrek/air
WORKDIR /app
ENTRYPOINT ["air"]