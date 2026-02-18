FROM alpine

RUN adduser -S -D -H -h /app appuser
WORKDIR /app

COPY ./clipr-linux-amd64 /app/clipr
RUN chmod +x /app/clipr

USER appuser

EXPOSE 8080
CMD [ "/app/clipr", "server" ]
