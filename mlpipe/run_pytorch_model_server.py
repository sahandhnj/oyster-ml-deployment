import time
import json
import yaml
import redis
import numpy as np
# import tensorflow as tf
# from keras.applications import imagenet_utils
from config.clistyle import bcolor
# from keras.models import model_from_json
from helpers import base64_decoding, NumpyEncoder
import torch

# Give along device details at input
device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")
# Also pass along grad_fn

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
graph_file = settings['model']['graph_file']
weights_file = settings['model']['weights_file']

model = None
use_gpu = False

def load_model(model_file_path, weights_file_path):
    global model

    model = torch.load(model_file_path)
    model_weights = torch.load(weights_file_path)
    model.load_state_dict(model_weights)
    # if use_gpu:
    #     model.cuda()
    print(bcolor.BOLD + "Loaded PyTorch model '{}' from disk and inserted weights from '{}'.".format(graph_file, weights_file) + bcolor.END)

    return model

def classify_process():
    model = load_model(model_dir + graph_file, model_dir + weights_file)

    while True:
        queue = rdb.lrange(
            settings['data_stream']['data_queue'], 0, settings['data_stream']['batch_size'] - 1)
        dataIDs = []
        batch = None

        for q in queue:
            q = q.decode("utf-8").replace("\'", "\"")
            q = json.loads(q)
            data = base64_decoding(q["data"], "float32", q["shape"])    # q["dtype"]
            # print("QSHAPE: ", q["shape"], q["filetype"]) # q["filetype"]
            if batch is None:
                batch = data
            else:
                batch = np.vstack([batch, data])
            dataIDs.append(q["id"])
            print(dataIDs)
            if len(dataIDs) > 0:
                print("Batch size: {}".format(batch.shape))
                predictions = model(torch.from_numpy(batch).to(device)) # Torch, set device

                print("PREDICTIONS: ", predictions)
                for (dataID, prediction) in zip(dataIDs, predictions):
                    output = []
                    
                    print("PREDICTION: ", prediction)

                    # Wrap in pytorch specific function later on
                    def torchTensor_detach_and_to_array(tensor):
                        with torch.no_grad():
                            tensor_detached = tensor.detach()
                            print("DETACHED: ", tensor_detached)
                            tensor_cpu = tensor_detached.cpu()
                            print("TENSOR CPU: ", tensor_cpu)
                            np_array = tensor_cpu.numpy()
                            print("NP ARR: ", np_array)
                            print(type(np_array))
                            print(np_array.dtype)
                            return np_array
                        
                    r = {"result": torchTensor_detach_and_to_array(prediction)}  # transform prediciton to np array
                    print(r)

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


# print(model)
# print(model.state_dict())