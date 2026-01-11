FROM node:20-alpine AS css
WORKDIR /app
COPY package.json tailwind.config.js ./
RUN npm install
COPY ui ./ui
RUN npm run build:css

FROM golang:1.24-alpine AS go
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
COPY --from=css /app/ui/static/css/output.css ./ui/static/css/output.css
RUN go build -o app ./cmd/web

FROM alpine:latest
WORKDIR /app
COPY --from=go /app/app .
COPY --from=go /app/ui ./ui
EXPOSE 4000
CMD ["./app"]