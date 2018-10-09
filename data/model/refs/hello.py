import os
import sys
import inspect
import uuid
import time
import json
import base64
import numpy
import numpy as np
from threading import Thread
import tensorflow as tf
from keras.models import model_from_json
from flask import Flask, request, jsonify
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


numpy.random.seed(42)
app = Flask(__name__)

# Flask variables
ALLOWED_EXTENSIONS = set(['txt', 'png', 'jpg', 'jpeg', 'wav'])

# Redis variables
db = redis.StrictRedis(host='localhost', port=6379, db=0)
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
    print("Loaded model from disk and inserted weights.")
    
    return loaded_model


# def base64_encoding(a):
#     return base64.b64encode(a).decode("utf-8")

# def base64_decoding(a, dtype):
#     shape = a.shape
#     if sys.version_info.major == 3:
#         a = bytes(a, encoding="utf-8")
#     a = np.frombuffer(base64.decodestring(a), dtype=dtype)
#     a = a.reshape(shape)
#     return a


# def allowed_file(filename):
#     return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

text = None

def afterwork(prediction):
    output = prediction[0][0]
    return output



def prep(user_input):
    return preprocessing(user_input)

def classify_process():
    model = load_model('model.json', 'model_weights.h5')
    # global text
    while True:
        queue = db.lrange(DATA_QUEUE, 0, BATCH_SIZE - 1)
        textIDs = []
        batch = None
        
        for q in queue:
            q = json.loads(q.decode("utf-8"))
            text = q["data"] # base64_decoding(
            print("TEXT1", text)
            if batch is None:
                batch = text
            else:
                batch = np.vstack([batch, text])
            
            textIDs.append(q["id"])
        # print("BATCH", batch)
        # batch2 = preprocessing(batch)
        # print("BATCH2", batch2)

        if len(textIDs) > 0:
            with graph.as_default():
                results = model.predict(batch) # predicitons

            results = afterwork(results)
            # for (textID, resultSet) in zip(textIDs, results):
            #     output = []                             
            #     r = {"ID": textID, "result": results}
            #     output.append(r)
                
            #     db.set(textID, json.dumps(output))
            for textID in textIDs:
                db.set(textID, json.dumps(float(results)))
            print(db)

            db.ltrim(DATA_QUEUE, len(textIDs), -1)

        time.sleep(SERVER_SLEEP)


@app.route('/predict', methods=["POST"])
def predict():
    
    data = {"success": False}

    if request.method == "POST":
        if request.files.get("text"):
            text = request.files["text"].read()
            # text = request.json["text"] # user_input
            
            print("TEXT2", text)
    
    # def prep(user_input):
    #     return preprocessing(user_input)
            text = text.copy(order="C")
            k = str(uuid.uuid4())
            d = {"id": k, "data": text} # base64_encoding(
            db.rpush(DATA_QUEUE, json.dumps(d))

            # print(d)
            while True:
                output = db.get(k)
                if output is not None:
                    output = output.decode("utf-8")
                    data["predictions"] = json.loads(output)
                    db.delete(k)
                    break
                time.sleep(CLIENT_SLEEP)
            data["success"] = True

    print(data)
    print(jsonify(data))
    return jsonify(data)



    # with graph.as_default():
    #     prediction = model.predict(preprocessing(user_input))

    # print("Input: %s. Prediction: %s" % (user_input, prediction))
    # return "Input: {}. Prediction: {}".format(user_input, prediction)


# @app.route("/")
# def hello():
#     return "Hello World!"

if __name__ == "__main__":
    print(("* Loading Keras model and Flask starting server..."
        "please wait until server has fully started"))
    print("* Starting model service...")
    t = Thread(target=classify_process, args=())
    t.deamon = True
    t.start()

    print("* Starting web service...")    
    app.run(host="0.0.0.0", port=int("5000"), debug=True)