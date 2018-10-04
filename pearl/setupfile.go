package pearl

import (
	"fmt"
	"log"
	"os"
)

func makeDockerFile() {
	writeToFiles(docker_file_static, "Dockerfile")
	writeToFiles(docker_compose_static, "docker-compose.yml")

}

func writeToFiles(content string, filename string) {
	fd, err := os.Create(filename)

	if err != nil {
		log.Fatal("Cannot create "+filename, err)
	}

	defer fd.Close()

	fmt.Fprintf(fd, content)
}

var docker_file_static = `FROM nvidia/cuda:9.0-cudnn7-devel

ENV MODELPATH /src
ENV USER sahand

RUN apt-get update && apt-get install -y software-properties-common

RUN add-apt-repository ppa:jonathonf/python-3.6

RUN apt-get update && apt-get install -y --no-install-recommends bzip2 python3.6 g++ \
    git graphviz libgl1-mesa-glx libhdf5-dev openmpi-bin wget unzip curl && \
    rm -rf /var/lib/apt/lists/*

RUN curl "https://bootstrap.pypa.io/get-pip.py" -o "get-pip.py"
RUN python3.6 get-pip.py

RUN useradd -m -s /bin/bash -N -u 1000 $USER && \
    mkdir -p $MODELPATH && \
    chown $USER $MODELPATH 

RUN pip install --upgrade pip

COPY /meta/requirements.txt $MODELPATH/requirements.txt

USER $USER
WORKDIR $MODELPATH 

RUN pip install -r requirements.txt --user


EXPOSE 5000
`

var docker_compose_static = `version: "2.4"
services:
  api:
    build:
     context: ./
    ports:
    - '5000'
    volumes:
        - /data/model:/data
    links:
      - redis
    command: python3.6 hello.py
  redis:
    image: redis:alpine
    volumes:
      - '/data/redis:/data'
    ports:
      - "6379"
`
