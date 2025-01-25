FROM golang:latest

COPY akademi /bin

CMD ["akademi", "daemon"]