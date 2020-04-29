FROM arm32v7/golang:latest

COPY . /src
WORKDIR /src

RUN apt-get update && apt-get install -y aria2 ffmpeg

RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

RUN go build -o main .

CMD ["./main"]