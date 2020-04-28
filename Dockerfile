FROM golang:latest

COPY . /src
WORKDIR /src
RUN apt-get update

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev aria2 \
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python \
  && pip3 install --upgrade pip

RUN pip3 install -r ./requirements.txt
RUN apt-get install -y ffmpeg
RUN go build -o main .

CMD ["./main"]