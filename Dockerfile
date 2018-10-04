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

USER $USER
WORKDIR $MODELPATH 

EXPOSE 5000

RUN pip install --user absl-py==0.5.0 astor==0.7.1 certifi==2018.8.24 chardet==3.0.4 Click==7.0 cycler==0.10.0 flasgger==0.9.1 Flask==1.0.2 gast==0.2.0 grpcio==1.15.0 h5py==2.8.0 idna==2.7 itsdangerous==0.24 Jinja2==2.10 jsonschema==2.6.0 Keras==2.2.2 Keras-Applications==1.0.6 Keras-Preprocessing==1.0.5 kiwisolver==1.0.1 Markdown==3.0.1 MarkupSafe==1.0 matplotlib==3.0.0 mistune==0.8.3 numpy==1.15.2 pandas==0.23.4 Pillow==5.2.0 protobuf==3.6.1 pyparsing==2.2.1 python-dateutil==2.7.3 pytz==2018.5 PyYAML==3.13 redis==2.10.6 requests==2.19.1 scikit-learn==0.20.0 scipy==1.1.0 seaborn==0.9.0 six==1.11.0 tensorboard==1.11.0 tensorflow==1.11.0 termcolor==1.1.0 urllib3==1.23 Werkzeug==0.14.1