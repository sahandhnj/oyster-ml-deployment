package pearl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func makeDockerFile() {
	docker_file_static = docker_file_static + "RUN pip install --user " + readRQFileIntoOneLine("data/model/requirements.txt")
	writeToFiles(docker_file_static, "Dockerfile")
	writeToFiles(docker_compose_static, "docker-compose.yml")
	writeToFiles("", ".dockerignore")

}

func writeToFiles(content string, filename string) {
	fd, err := os.Create(filename)

	if err != nil {
		log.Fatal("Cannot create "+filename, err)
	}

	defer fd.Close()

	fmt.Fprintf(fd, content)
}

func readRQFileIntoOneLine(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var stringArray []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stringArray = append(stringArray, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return strings.Join(stringArray, " ")
}

var docker_file_static = `FROM python:alpine3.7

ENV MODELPATH /src

RUN mkdir -p $MODELPATH
WORKDIR $MODELPATH 

EXPOSE 5000
`
var docker_file_static_ubuntu = `FROM ubuntu:18.04

ENV MODELPATH /src
ENV USER sahand

RUN apt-get update && apt-get install -y software-properties-common

RUN add-apt-repository ppa:jonathonf/python-3.6

RUN apt-get update && apt-get install -y --no-install-recommends bzip2 python3.6 g++ \
	python3-distutils git graphviz libgl1-mesa-glx libhdf5-dev openmpi-bin wget unzip curl && \
    rm -rf /var/lib/apt/lists/*

RUN curl "https://bootstrap.pypa.io/get-pip.py" -o "get-pip.py"
RUN python3.6 get-pip.py

RUN useradd -m -s /bin/bash -N -u 1000 $USER && \
    mkdir -p $MODELPATH && \
    chown $USER $MODELPATH 

RUN pip install --upgrade pip

USER $USER
WORKDIR $MODELPATH 

EXPOSE 5000
`

var docker_file_static_gpu = `FROM nvidia/cuda:9.0-cudnn7-devel

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

USER $USER
WORKDIR $MODELPATH 

EXPOSE 5000
`
var docker_compose_static = `version: "2.4"
services:
  api:
    build:
     context: ./
    ports:
    - '5000:5000'
    volumes:
        - './data/model:/src'
    command: python hello.py
`

var docker_compose_static_redis = `version: "2.4"
services:
  api:
    build:
     context: ./
    ports:
    - '5000:5000'
    volumes:
        - './data/model:/src'
    links:
      - redis
    command: python3.6 --version
  redis:
    image: redis:alpine
    volumes:
      - './data/redis:/data'
    ports:
      - "6379"
`
