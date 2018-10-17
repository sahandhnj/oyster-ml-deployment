from flask import Flask, request, jsonify
import torch
import numpy as np
import torchvision
from torchvision import datasets, models, transforms
import matplotlib.pyplot as plt
from werkzeug.utils import secure_filename
import os
from PIL import Image
from helpers import NumpyEncoder
import json
import yaml
import uuid
import redis
import time

from helpers import base64_encoding, get_dtype

with open("./config/settings.yaml", 'r') as stream:
    try:
        settings = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)

# TBD after merge
with open("./config/allowedExtns.yaml", 'r') as stream:
    try:
        allowed_extensions = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)




app = Flask(__name__)
rdb = redis.StrictRedis(
    host=settings['redis']['host'],
    port=settings['redis']['port'],
    db=settings['redis']['db']
)

# rdb.flushall()

model = None
use_gpu = False

def load_model(model_file_path, weights_file_path):
    global model

    model = torch.load(model_file_path)
    model_weights = torch.load(weights_file_path)
    model.load_state_dict(model_weights)
    # if use_gpu:
    #     model.cuda()

    return model


def torchTensor_to_npArray(tensor):
    """Detatch Tensor, write to cpu, conver to numpy array
    """
    # Do checks first
    with torch.no_grad():
        print("TENS1: ", tensor)
        tensor_detatched = tensor.detach()
        print("TENS2: ", tensor_detatched)
        tensor_cpu = tensor_detatched.cpu()
        print("TENS3: ", tensor_cpu)
        np_array = tensor_cpu.numpy()
        print("ARR: ", np_array)
    return np_array

def preprocessing(image):
    # Define preprocessing transformation
    composed = transforms.Compose([
        transforms.Resize(256),
        transforms.CenterCrop(224),
        transforms.ToTensor(),
        transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
    ])
    preprocessed = composed(image)
    preprocessed.unsqueeze_(0)  # Set correct batch dimension

    return preprocessed

device = torch.device("cuda:0" if torch.cuda.is_available() else "cpu")

model = load_model('./runs/pytorch/model/resnet18_sinp.json', './runs/pytorch/model/resnet18_sinp_weights.h5')


def classify_process(model, inputs):
    model.eval()
    model.training = False
    output = model(inputs)
    return output


def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in set(settings['data_stream']['allowed_extensions'])

def get_file_type(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower()


@app.route("/predict", methods=["POST"])
def predict():
    
    data = {"success": False}

    if request.method == "POST":
        # Check if file in inputted
        if 'data' not in request.files:
            flash("Nof ile part")
            raise ValueError("No file part")
        # print("METHOD: ", request.method)
        file = request.files['data']    # Redundant?
        filetype = get_file_type(file.filename)
        # print("FILE: ", file)
        if file.filename == '':
            flash("No selected file")
            raise ValueError("No selected file")
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            if request.files.get('data'):
                user_input = request.files["data"]  #.read()
                # print("UI1: ", user_input)
                # user_input = request.files["data"].read()
                # print("UI1R: ", user_input)
                if (filetype in ['jpg', 'jpeg', 'png']):
                    user_input = Image.open(user_input)
                else:
                    pass
        
                preprocessed_input = preprocessing(user_input)
                # print("PEPROCESSED: ", preprocessed_input)
                
                # Get file properties
                if filetype in ['jpg', 'jpeg', 'png']:
                    fileshape = np.array(preprocessed_input).shape
                else:
                    fileshape = preprocessed_input.shape
                print("FILESHAPE: ", fileshape)

                array_dtype = get_dtype(preprocessed_input)
                # HANDLE TORCH TENSORS
                # print("DTYPE: ", array_dtype)
                # if array_dtype == 'torch._':
                #     print("TORCH!")
                #     ### convert torch to numpy
                #     output_np = torchTensor_to_npArray(output1)
                
                preprocessed_input = preprocessed_input.numpy()
                # print("NP PREP: ", preprocessed_input)               
                preprocessed_input = preprocessed_input.copy(order="C")
                encoded_input = base64_encoding(preprocessed_input)

                k = str(uuid.uuid4())
                d = {
                    "id": k,
                    "filename": filename,
                    "filetype": filetype,
                    "shape": fileshape,
                    "dtype": array_dtype,
                    "data": encoded_input
                }
                rdb.rpush(settings['data_stream']['data_queue'], json.dumps(d))

                # Can i also send and receive torch tensors via redis?
                # or setup automated function to resore them from np array

                while True:
                    output = rdb.get(k)
                    if output is not None:
                        output = output.decode("utf-8")
                        data["summary"] = json.loads(output)
                        rdb.delete(k)
                        break
                    
                    time.sleep(settings['data_stream']['client_sleep'])
                data["success"] = True
                # print("OUTPUT: ", output)
    return jsonify(data)

        # print("USER_INPUT: ", user_input)
        
        
        
    #     # print("DEVICE: ", device)
    #     image = Image.open(file)
    #     # print("IMAGE: ", image)
    #     tfd1 = preprocessing(image)
    #     # print("PREPPED:", tfd1)
    #     output1 = classify_process(model=model, inputs=tfd1.to(device))
    #     # print("OUTPUT: ", output1)
    #     ### convert torch to numpy
    #     output_np = torchTensor_to_npArray(output1)
    #     # print("OUTPUT NP: ", output_np)
    #     output_json = json.dumps(output_np, cls=NumpyEncoder)
    #     # print("OUTPUT JSON", output_json)

    # return output_json


if __name__ == "__main__":
    app.run(port=5001)