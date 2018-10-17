import os
import requests
import numpy as np
from keras.applications import imagenet_utils

cwd = os.getcwd()
image_dir = cwd + "/data/images/"

image_files = os.listdir(image_dir)

def stream(file):
    # load the input image and construct the payload for the request
    filecontent = open(image_dir + file, "rb") #.read()

    # submit the request
    r = requests.post("http://localhost:5000/predict", files={"data": filecontent})

    return r.json()


for i in image_files:
    resp = stream(i)
    prediction = resp["summary"][0]['result']
    np_prediction = np.array([prediction])
    pred = imagenet_utils.decode_predictions(np_prediction)
    pred_name = pred[0][0][1]
    pred_val = pred[0][0][2]

    print("PREDICTED: {0:<30s} \tCERTAINTY: {1:.3f}".format(pred_name, pred_val))


