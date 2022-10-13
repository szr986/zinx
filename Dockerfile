FROM golang:alpine


WORKDIR /build

COPY . .

RUN go build . 

EXPOSE 8888

CMD ["/build/app"]