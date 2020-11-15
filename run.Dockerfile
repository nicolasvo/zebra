FROM alpine:latest
WORKDIR /app
COPY . .
CMD ["sh", "run.sh"]