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
# print(os.listdir(data_dir))
# imgpath1 = data_dir + '0013035.jpg'
# imgpath2 = data_dir + '16838648_415acd9e3f.jpg'
 # img1 = Image.open(open(imgpath1, 'rb'))
def stream(file):
    filecontent = open(file, 'rb')
    r = requests.post("http://localhost:5000/predict", files={"data": filecontent})
    return r.json()
 # resp = stream(imgpath1)
# # resp = stream(imgpath2)
# print(resp)
class_names = ['ants', 'bees']

def process_output(outputs, image):
    print("OUTPUTS: ", outputs)
    outputs = torch.Tensor(outputs)
    print("OUTPUTS TENS: ", outputs)
    _, preds = torch.max(outputs, 1)
    print("PREDS: ", preds)
    for j in range(outputs.size()[0]):
        ax = plt.subplot()
        ax.axis('off')
        ax.set_title('predicted: {}'.format(class_names[preds[j]]))
        plt.imshow(image)
        plt.show()

if __name__ == "__main__":
    # process_output(resp, image=img1)
    image_files = os.listdir(data_dir)
    for i in image_files:
        resp = stream(data_dir + i)
        img = Image.open(open(data_dir + i, 'rb'))
        process_output(resp, image=img) 