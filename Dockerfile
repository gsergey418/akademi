FROM golang:latest

COPY akademi /bin

CMD ["akademi", "daemon", "--rpc-addr", "0.0.0.0:3855"]