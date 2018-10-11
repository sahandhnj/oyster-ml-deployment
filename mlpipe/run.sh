#!/bin/bash
python run_web_server.py & 2>&1 | tee output
python run_model_server.py  2>&1 | tee output
