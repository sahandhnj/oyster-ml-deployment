from flask import Flask, request, jsonify
import torch
import numpy as np
import torchvision
from torchvision import datasets, models, transforms
import matplotlib.pyplot as plt
import os
from PIL import Image
from helpers import NumpyEncoder
import json

app = Flask(__name__)

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


@app.route("/predict", methods=["POST"])
def predict():
    
    if request.method == "POST":
        print("METHOD: ", request.method)
        file = request.files['data']
        print("FILE: ", file)
        print("DEVICE: ", device)
        image = Image.open(file)
        print("IMAGE: ", image)
        tfd1 = preprocessing(image)
        print("PREPPED:", tfd1)
        output1 = classify_process(model=model, inputs=tfd1.to(device))
        print("OUTPUT: ", output1)
        ### convert torch to numpy
        output_np = torchTensor_to_npArray(output1)
        print("OUTPUT NP: ", output_np)
        output_json = json.dumps(output_np, cls=NumpyEncoder)
        print("OUTPUT JSON", output_json)

    return output_json


if __name__ == "__main__":
    app.run()