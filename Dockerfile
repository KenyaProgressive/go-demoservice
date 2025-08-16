FROM golang:alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
EXPOSE 8000
ENTRYPOINT ["go", "run", "."]
