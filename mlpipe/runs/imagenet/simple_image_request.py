import requests

import os

cwd = os.getcwd()
image_dir = cwd + "/data/images/"

image_files = os.listdir(image_dir)
image1 = image_files[0]
print(image1)

def stream(file):
    # load the input image and construct the payload for the request
    filecontent = open(image_dir + file, "rb") #.read()

    # submit the request
    r = requests.post("http://localhost:5000/predict", files={"data": filecontent})

    return r.json()


for i in image_files:
    resp = stream(i)
    prediction = resp["summary"][0]['result'][1]
    pred_val = resp["summary"][0]['result'][2]
    print(prediction, pred_val)



