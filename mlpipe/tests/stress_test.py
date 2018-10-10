from threading import Thread
import requests
import time
import requests
import json

KERAS_REST_API_URL = "http://localhost:5000/predict"
# CURL_REQUEST = "curl --header \"Content-Type: application/json\" --request POST --data '{\"text\": \"you like the movie\"}' http://localhost:5000/predict"

payload = {"text": "you like the movie"}

NUM_REQUESTS = 300
SLEEP_COUNT = 0.05

# r = requests.post(KERAS_REST_API_URL, data=json.dumps(payload), headers=headers)
# resp = r.json()

headers = {"Content-Type": "application/json"}
def call_predict_endpoint(n):


    text = "you like the movie"
    payload = {"text": text}
    r = requests.post(KERAS_REST_API_URL, data=json.dumps(payload), headers=headers).json()

    if r["success"]:
        print("[INFO] thread {} OK".format(n))
    else:
        print("[INFO] thread {} FAILED".format(n))
    print("Success: True ")

for i in range(0, NUM_REQUESTS):
    t = Thread(target=call_predict_endpoint, args=(i,))
    t.daemon = True
    t.start()
    time.sleep(SLEEP_COUNT)

time.sleep(300)



# print("Success: ", resp["success"])
# print("Result: ", resp['predictions'][0]['result'])
#


#
# def call_predict_endpoint(n):
#     curl = CURL_REQUEST
#     payload = {"data": "you like the movie"}
#     r = requests.post(KERAS_REST_API_URL, data=payload)
