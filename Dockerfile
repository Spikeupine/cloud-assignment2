FROM golang:1.22.2

LABEL developer="Solveig Langbakk, Sofia Mikkelsen, Trygve Sollund"
LABEL maintainer="trygve.sollund@ntnu.no"
LABEL  stage=builder

WORKDIR /go/src/app/

#Copy Directory Structure of the project
COPY ./cmd /go/src/app/cmd
COPY ./database /go/src/app/database
COPY ./external /go/src/app/external
COPY ./internal /go/src/app/internal
#Download go mod dependencies
COPY ./go.mod  /go/src/app
COPY ./go.sum /go/src/app
RUN go mod download && go mod verify
#
RUN CGO_ENABLED=0 GOOS=linux go build  -a -ldflags '-extldflags "-static"'   -o ./run ./cmd

EXPOSE 8080
CMD ["./run"]