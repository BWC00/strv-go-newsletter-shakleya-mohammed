# Build environment
# -----------------
FROM golang:1.20-alpine as build-env
WORKDIR /strv-go-newsletter-shakleya-mohammed

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/api \
    && go build -ldflags '-w -s' -a -o ./bin/migrate ./cmd/migrate


# Deployment environment
# ----------------------
FROM alpine

COPY --from=build-env /strv-go-newsletter-shakleya-mohammed/bin/api /strv-go-newsletter-shakleya-mohammed/
COPY --from=build-env /strv-go-newsletter-shakleya-mohammed/bin/migrate /strv-go-newsletter-shakleya-mohammed/
COPY --from=build-env /strv-go-newsletter-shakleya-mohammed/migrations /strv-go-newsletter-shakleya-mohammed/migrations

EXPOSE 8080
CMD ["/strv-go-newsletter-shakleya-mohammed/api"]