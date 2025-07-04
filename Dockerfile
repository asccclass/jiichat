FROM alpine

WORKDIR /app
COPY ./app /app
ENTRYPOINT ["/app/app"]
