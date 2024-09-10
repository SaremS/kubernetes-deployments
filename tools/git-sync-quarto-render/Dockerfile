FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o quarto-render-job .

FROM registry.k8s.io/git-sync/git-sync:v4.2.4
COPY --from=builder /app/quarto-render-job /quarto-render-job
USER root
RUN chmod 755 /quarto-render-job
RUN chown 65533:65533 /quarto-render-job
USER 65533:65533
