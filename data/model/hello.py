import tensorflow as tf
import numpy
from keras.preprocessing import sequence
from keras.models import model_from_json
import keras
from flask import Flask, request
import redis
import inspect
from keras.models import model_from_json
import numpy
import os

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

# Redis variables
db = redis.StrictRedis(host='localhost', port=6397, db=0)
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


model = load_model('model.json', 'model_weights.h5')

@app.route('/predict', methods=['POST'])
def predict():
    user_input = request.json['text']

    with graph.as_default():
        prediction = model.predict(preprocessing(user_input))

    print("Input: %s. Prediction: %s" % (user_input, prediction))
    return "Input: {}. Prediction: {}".format(user_input, prediction)


# @app.route("/")
# def hello():
#     return "Hello World!"

if __name__ == "__main__":
    print(("* Loading Keras model and Flask starting server..."
        "please wait until server has fully started"))
    app.run(host="0.0.0.0", port=int("5000"), debug=True)