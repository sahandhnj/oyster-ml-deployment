import time
import json
import yaml
import redis
import numpy as np
import tensorflow as tf


from config.clistyle import bcolor
from keras.models import model_from_json
from helpers import base64_decoding

with open("./config/settings.yaml", 'r') as stream:
    try:
        settings = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)


rdb = redis.StrictRedis(
    host=settings['redis']['host'],
    port=settings['redis']['port'],
    db=settings['redis']['db']
)
model_dir = settings['model']['pathdir']


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


def classify_process():
    model = load_model(model_dir + 'model.json', model_dir + 'model_weights.h5')

    while True:
        queue = rdb.lrange(
            settings['data_stream']['data_queue'], 0, settings['data_stream']['batch_size'] - 1)
        dataIDs = []
        batch = None
        
        for q in queue:
            q = q.decode("utf-8").replace("\'", "\"")
            q = json.loads(q)  ###
            data = base64_decoding(q["data"], q["dtype"], q["shape"])  ###
            if batch is None:
                batch = data
            else:
                batch = np.vstack([batch, data])  # if already data in queue add a new layer
            dataIDs.append(q["id"])
            # Check if it fits in batch and processing is needed
            if len(dataIDs) > 0:
                print("Batch size: {}".format(batch.shape))

                with graph.as_default():
                    predictions = model.predict(batch)

                for (dataID, predictionSet) in zip(dataIDs, predictions):
                    output = []
                    for prediction in predictionSet:
                        r = {"result": float(
                            prediction)}  # modify prediction as non-array so it can be stored to redis db
                        output.append(r)
                    output.append({"input": {"uid": dataID, "dtype": q["dtype"], "shape": batch.shape}})

                    rdb.set(dataID, json.dumps(output))
                rdb.ltrim(settings['data_stream']['data_queue'], len(dataIDs), -1)
            time.sleep(settings['data_stream']['server_sleep'])


if __name__ == "__main__":
    classify_process()