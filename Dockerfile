FROM golang:1.17.7
COPY ./ /home/services/tg-users-service
WORKDIR /home/services/tg-users-service/cmd
RUN go mod tidy
ENTRYPOINT ["go","run","cmd.go"]