#!/bin/bash
#!bin/bash
if [ ! -d "keraspipeline" ]; then
  tar -xzvf model.tar.gz
fi

cd keraspipeline/servers
python3.6 run_keras_web_server.py & 2>&1 | tee output
python3.6 run_keras_model_server.py  2>&1 | tee output
