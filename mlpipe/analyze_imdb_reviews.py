import os
import time
import requests
import json
import yaml
with open("./config/settings.yaml", 'r') as stream:
    try:
        settings = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)


cwd = os.getcwd()
data_dir = cwd + "/data/tweets/"

tweet_files = os.listdir(data_dir)
tweet = tweet_files[0]
print(tweet)

def feedstream(file, sleep=0.05, verbose=1, *args, **kwargs):
    API_ENDPOINT = settings['model']['api_endpoint']
    headers = None
    filecontent = open(data_dir + file)
    # payload = {"data": open("tweet.csv")}

    r = requests.post(API_ENDPOINT, files={"data": filecontent})
    if verbose == 1:
        if r.json()["success"]:
            print("[INFO] feeding {}".format(file))
        else:
            print("[INFO] feed for {} failed".format(file))
    else:
        pass
    # time.sleep(sleep)    
    
    return r.json()



# Retrieve
results = []
for t in tweet_files:
    resp = feedstream(t, verbose=0)
    prediction = resp['summary'][0]['result']
    success = resp["success"]
    inputfile = resp['summary'][1]['input']['filename']
    filecontent = open(data_dir + t).read()
    print("REVIEW: {0:<30s} \t SENTIMENT: {1:.3f} \t (file={2}, analysis={3})".format(filecontent, prediction, inputfile, success))
   
    results.append(resp)
        


# print(results)

# todo: 
# merge with stress test, add threading
# add autoselect function from feedstream
# add statistics
# add writing results to csv
# diversify feedstream with images and text options
# Test feedstream with images 


# COMMENT
# https://stackoverflow.com/questions/29526688/send-a-csv-file-using-requests-2-2-1-in-python-2-7-6