FROM golang:1.15 as builder

COPY . /source
WORKDIR /source
# build statically linked binary
RUN env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a

FROM busybox:stable
RUN mkdir /honeypot
COPY --from=builder /source/ssh_honeypot /honeypot/
RUN chmod +x /honeypot/ssh_honeypot

ENV SSH_PORT=22
ENV WAIT_DURATION=15
ENV MAX_ATTEMPTS=20
ENV LOG_FILE=/honeypot/attempts.log

ENTRYPOINT /honeypot/ssh_honeypot -p $SSH_PORT -w $WAIT_DURATION -l $LOG_FILE -c $MAX_ATTEMPTS

