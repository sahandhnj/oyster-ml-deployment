import requests
import json

KERAS_REST_API_URL = "http://localhost:5000/predict"
# CURL_REQUEST = "curl --header \"Content-Type: application/json\" --request POST --data '{\"text\": \"you like the movie\"}' http://localhost:5000/predict"

headers = {"Content-Type": "application/json"}
payload = {"text": "you like the movie"}

r = requests.post(KERAS_REST_API_URL, data=json.dumps(payload), headers=headers)
resp = r.json()

print("Success: ", resp["success"])
print("Result: ", resp['predictions'][0]['result'])


# print(resp['predictions'][0]['input'])


