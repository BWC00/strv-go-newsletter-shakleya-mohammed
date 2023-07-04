FROM golang:1.20-alpine
WORKDIR /strv-go-newsletter-shakleya-mohammed

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate

CMD ["/strv-go-newsletter-shakleya-mohammed/bin/api"]
EXPOSE 8080