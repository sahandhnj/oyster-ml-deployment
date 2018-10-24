import os, sys, inspect, time, json, yaml
import redis
import numpy as np
import tensorflow as tf
from keras.applications import imagenet_utils
# sys.path.insert(1, os.path.join(sys.path[0], '../..'))  # insert mlpipe root to path
mlpipe_root = os.path.abspath("..")
sys.path.insert(0, mlpipe_root)

from config.clistyle import bcolor
from keras.models import model_from_json
from servers.helperfunctions import base64_decoding, NumpyEncoder

with open(mlpipe_root + "/settings.yaml", 'r') as stream:
    try:
        settings = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)


rdb = redis.StrictRedis(
    host=settings['redis']['host'],
    port=settings['redis']['port'],
    db=settings['redis']['db']
)


def get_paths(root_path):
    model_dir = root_path + settings['model']['pathdir']
    graph_file = settings['model']['graph_file']
    weights_file = settings['model']['weights_file']
    graph = model_dir + graph_file
    weights = model_dir + weights_file
    
    return graph, graph_file, weights, weights_file


def load_model(model_file_path, weights_file_path, graph_file, weights_file):
    global model

    with open("{}".format(model_file_path), 'r') as model_json_file:
        loaded_model_json = model_json_file.read()
    loaded_model = model_from_json(loaded_model_json)
    loaded_model.load_weights("{}".format(weights_file_path))

    global graph
    graph = tf.get_default_graph()
    print(bcolor.BOLD + "Loaded model '{}' from disk and inserted weights from '{}'.".format(graph_file, weights_file) + bcolor.END)

    return loaded_model


def classify_process():
    
    graph_path, graph_file, weights_path, weights_file = get_paths(mlpipe_root)
    model = load_model(graph_path, weights_path, graph_file, weights_file)

    while True:
        queue = rdb.lrange(
            settings['data_stream']['data_queue'], 0, settings['data_stream']['batch_size'] - 1)
        dataIDs = []
        batch = None
        
        for q in queue:
            q = q.decode("utf-8").replace("\'", "\"")
            q = json.loads(q)
            data = base64_decoding(q["data"], q["dtype"], q["shape"])
            print("QSHAPE: ", q['shape'], q["filetype"])
            
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
                # print("PREDICITONS: ", predictions, type(predictions))
                
                # This if statement possibly move out of model server, since imagenet specific
                # if (q["filetype"] in ['jpg', 'jpeg', 'png']):
                #     predictions = imagenet_utils.decode_predictions(predictions)
                # else:
                #     pass
                # print("PREDICTIONS: ", predictions)

                for (dataID, prediction) in zip(dataIDs, predictions):
                    output = []
                    # for prediction in predictionSet:
                    # print("PREDICTION: ", prediction, type(predictions))
                    r = {"result": prediction}  # float() modify prediction as non-array so it can be stored to redis db
                    output.append(r)
                    output.append({
                        "input": {
                            "uid": dataID,
                            "filename": q["filename"],
                            "filetype": q["filetype"],
                            "dtype": q["dtype"],
                            "shape": batch.shape
                            }
                        })

                    rdb.set(dataID, json.dumps(output, cls=NumpyEncoder))
                rdb.ltrim(settings['data_stream']['data_queue'], len(dataIDs), -1)
            time.sleep(settings['data_stream']['server_sleep'])


if __name__ == "__main__":
    classify_process()
