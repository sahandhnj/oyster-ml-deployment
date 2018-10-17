import os
import requests
import numpy as np
import torch
import torchvision
from torchvision import datasets, models, transforms
from PIL import Image
import matplotlib.pyplot as plt

# Define path
data_dir = '/home/bloks/Projects/Sentriq/apiclient/mlpipe/runs/pytorch/data/hymenoptera_data/test/'

def stream(file):
    filecontent = open(file, 'rb')

    r = requests.post("http://localhost:5000/predict", files={"data": filecontent})

    return r.json()

class_names = ['ants', 'bees']

def process_output(outputs, image):
    # outputs = torch.Tensor(outputs)
    x = torch.max(outputs, 0)
    print(x)
    
    _, preds = torch.max(outputs, 1)
    # print("PREDS: ", preds)
    for j in range(outputs.size()[0]):
        ax = plt.subplot()
        ax.axis('off')
        ax.set_title('predicted: {}'.format(class_names[preds[j]]))
        plt.imshow(image)
        plt.show()


if __name__ == "__main__":

    image_files = os.listdir(data_dir)

    for i in image_files:
        resp = stream(data_dir + i)
        # print(resp)
        res_pred = resp['summary'][0]['result']
        print(res_pred)
        prediction = torch.Tensor([res_pred])
        print(prediction)
        print(prediction.size())
        img = Image.open(open(data_dir + i, 'rb'))
        process_output(prediction, image=img)