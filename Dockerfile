FROM ubuntu:18.04
ENV MODELPATH /src

RUN apt-get update && apt-get install -y software-properties-common

RUN add-apt-repository ppa:jonathonf/python-3.6 
RUN apt-get update && apt-get install -y python3.6 curl python-pip python-dev build-essential 
	
RUN python3.6 --version
RUN pip --version

RUN pip install --upgrade pip

WORKDIR $MODELPATH 
EXPOSE 5000
RUN pip install --user Flask==1.0.2 h5py==2.8.0 Keras==2.2.2 tensorflow==1.11.0 redis==2.10.6