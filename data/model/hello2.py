import os
import time
import uuid
import json
import base64
import inspect
import sys
from threading import Thread
from clistyle import bcolor
import numpy
import numpy as np
import tensorflow as tf
from keras.models import model_from_json
from flask import Flask, request, jsonify
import requests
import redis
try:
    import preprocessing as prepmod
    if hasattr(prepmod, 'preprocessing') and inspect.isfunction(prepmod.preprocessing):
        from preprocessing import preprocessing
        print("Preprocessing file available and loaded into vessel.")
    else:
        raise TypeError("Preprocessing file inserted, but does not contain function called 'preprocessing'.")
except (ImportError):
    print("No preprocessing file inserted.")


from io import BytesIO

numpy.random.seed(42)
app = Flask(__name__)

# Flask variables
# ALLOWED_EXTENSIONS = set(['txt', 'png', 'jpg', 'jpeg', 'wav'])

# Redis variables
rdb = redis.StrictRedis(host='redis', port=6379, db=0)
DATA_QUEUE = "data_queue"
BATCH_SIZE = 32
SERVER_SLEEP = 0.25
CLIENT_SLEEP = 0.25


def load_model(model_file_path, weights_file_path):
    global model
    with open("{}".format(model_file_path), 'r') as model_json_file:
        loaded_model_json = model_json_file.read()
    loaded_model = model_from_json(loaded_model_json)
    loaded_model.load_weights("{}".format(weights_file_path))
    
    global graph
    graph = tf.get_default_graph()
    print(bcolor.BOLD + "Loaded model from disk and inserted weights." + bcolor.END)

    return loaded_model

def base64_encoding(array):
    return base64.b64encode(array).decode("utf-8")

def base64_decoding(array, dtype, shape):
    if sys.version_info.major ==3:
        array = bytes(array, encoding="utf-8")
    array = np.frombuffer(base64.decodestring(array), dtype=dtype)
    array = array.reshape(shape)
    return array

IMPUT_SIZE = (1, 500)

def classify_process():
    model = load_model('model.json', 'model_weights.h5')

    while True:
        queue = rdb.lrange(DATA_QUEUE, 0, BATCH_SIZE -1)
        dataIDs = []
        batch = None
        
        for q in queue:
            q = q.decode("utf-8").replace("\'", "\"")
            q = json.loads(q)
            # print(r["data"])
            data = base64_decoding(q["data"], 'float32', (1,500))
            if batch is None:
                batch = data
            else:
                batch = np.vstack([batch, data]) # if already data in queue add a new layer
            dataIDs.append(q["id"])
            # Check if it fits in batch and processing is needed
            if len(dataIDs) > 0:
                print("Batch size: {}".format(batch.shape))
                with graph.as_default():
                    predictions = model.predict(batch)               
            
                for (dataID, prediction) in zip(dataIDs, predictions):            
                    output = []
                    r = {"id": dataID, "result": float(prediction)} # modify prediction as non-array so it can be stored to redis db
                    output.append(r)
                    # print(output)

                rdb.set(dataID, json.dumps(output))
            rdb.ltrim(DATA_QUEUE, len(dataIDs), -1)
        time.sleep(SERVER_SLEEP)


@app.route('/predict', methods=["POST"])
def predict():

    data = {"success": False}

    if request.method == "POST":
        user_input = request.json["text"]      
        preprocessed_input = preprocessing(user_input)
        encoded_input = base64_encoding(preprocessed_input)
        # MAKE C-CONTIGOUS?????
        
        k = str(uuid.uuid4())
        d = {"id": k, "shape": preprocessed_input.shape, "data": encoded_input}
        rdb.rpush(DATA_QUEUE, json.dumps(d))    # dump the preprocessed input as a numpy array

        while True:
            output = rdb.get(k)
            # print(output)

            if output is not None:
                output = output.decode("utf-8")
                data["predictions"] = json.loads(output)

                rdb.delete(k)
                break
            
            time.sleep(CLIENT_SLEEP)
        data["success"] = True
    
    return jsonify(data)    


# @app.route("/")
# def hello():
#     return "Hello World!"

if __name__ == "__main__":
    print(("* Loading Keras model and Flask starting server..."
        "please wait until server has fully started"))
    
    print("* Starting model service... ")
    t = Thread(target=classify_process, args=())
    t.daemon
    t.start()
    
    
    print("* Starting web service...")
    app.run(host="0.0.0.0", port=int("5000"), debug=True)