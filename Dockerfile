FROM golang:1.22 

RUN go version 
ENV GOPATH=/ 

COPY ./ ./ 

RUN go mod download 
RUN go build -o cad-keeper-auth ./cmd/main.go 

CMD ["./cad-keeper-auth"]