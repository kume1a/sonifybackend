FROM golang:1.24.1

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download

COPY *.go ./
COPY ./internal ./internal
COPY ./sql ./sql

RUN CGO_ENABLED=0 GOOS=linux go build -o /sonifybin

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Install dependencies
RUN mkdir ~/.local
RUN mkdir ~/.local/bin
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o ~/.local/bin/yt-dlp
RUN chmod a+rx ~/.local/bin/yt-dlp

RUN apt update
RUN apt-get -y install make
RUN apt install ffmpeg -y

EXPOSE 8000

ENV SONIFY_ENV=production
ENV PATH="$PATH:/root/.local/bin"

# ENTRYPOINT ["tail", "-f", "/dev/null"]

CMD ["/bin/bash", "-c", "make migrate-prod;/sonifybin"]