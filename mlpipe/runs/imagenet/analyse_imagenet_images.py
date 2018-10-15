import os
import time
import requests
import json
import yaml

cwd = os.getcwd()
print(cwd)


with open("../../config/settings.yaml", 'r') as stream: # "../../config/settings.yaml"
    try:
        settings = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)

cwd = os.getcwd()
data_dir = cwd + "/data/images/"

image_files = os.listdir(data_dir)
image = image_files[0]
print(image)


def feedstream(file, sleep=0.05, verbose=1, *args, **kwargs):
    API_ENDPOINT = settings['model']['api_endpoint']
    headers = None
    filecontent = open(data_dir + file)
    # payload = {"data": open("tweet.csv")}

    r = requests.post(API_ENDPOINT, files={"data": filecontent}).json()
    if verbose == 1:
        if r.json()["success"]:
            print("[INFO] feeding {}".format(file))
        else:
            print("[INFO] feed for {} failed".format(file))
    else:
        pass
    # time.sleep(sleep)

    return r.json()

results = []
for image in image_files:
    resp = feedstream(image, verbose=0)
    prediction = resp["summary"]
    print(prediction)