FROM golang:1.22.0 AS build-stage

WORKDIR /app

COPY go.mod go.sum Makefile ./.env.production ./
RUN go mod download

COPY *.go ./
COPY ./internal ./internal
COPY ./sql ./sql

RUN CGO_ENABLED=0 GOOS=linux go build -o /sonifybin

## Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

## Setup sudo
# RUN apt update
# RUN apt install sudo

# RUN adduser --disabled-password --gecos '' admin
# RUN adduser admin sudo
# RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# USER admin

## Install packages
# RUN sudo apt -y install software-properties-common

RUN apt update
RUN apt install yt-dlp -y
RUN apt-get -y install make

EXPOSE 8000

ENV SONIFY_ENV=production

# ENTRYPOINT ["tail", "-f", "/dev/null"]

CMD ["/bin/bash", "-c", "make migrate-prod;/sonifybin"]