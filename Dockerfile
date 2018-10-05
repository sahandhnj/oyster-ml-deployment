FROM python:alpine3.7

ENV MODELPATH /src

RUN mkdir -p $MODELPATH
WORKDIR $MODELPATH 

EXPOSE 5000
RUN pip install --user Flask==1.0.2