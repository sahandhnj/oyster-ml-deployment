FROM nvidia/cuda:9.0-cudnn7-devel

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
