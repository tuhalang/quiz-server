FROM golang

ENV GO111MODULE on

RUN go version

COPY . ~/src

WORKDIR ~/src

RUN go mod download

RUN go build -o quiz

RUN chmod +x quiz

EXPOSE 8080

ENTRYPOINT [ "./quiz" ]