FROM alpine

RUN adduser -S -D -H -h /app appuser
WORKDIR /app

COPY ./clippy-linux-amd64 /app/clippy
RUN chmod +x /app/clippy

USER appuser

EXPOSE 8080
CMD [ "/app/clippy", "server" ]
