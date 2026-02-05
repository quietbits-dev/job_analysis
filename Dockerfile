FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -buildvcs=false -o main .

RUN CGO_ENABLED=0 GOOS=linux go build -o crawler main.go


FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont


ENV CHROME_BIN=/usr/bin/chromium-browser

COPY --from=builder /app .
COPY --from=builder /app/parser/gethtml.go .
COPY --from=builder /app/parser/getlink.go .
COPY --from=builder /app/repo/insert.go .


CMD ["./crawler"]