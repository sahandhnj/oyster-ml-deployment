"""
Simple script to obtained trained models from keras.
Model folder is in .gitignore due to the large size of the models.
"""

import os
from keras.applications import ResNet50

model = ResNet50(weights="imagenet")
model.trainable = False

cwd = os.path.abspath('.')
model_dir = cwd + "/model/" 

model_json = model.to_json()
with open(model_dir + "resnet50_imagenet.json", "w") as json_file:
    json_file.write(model_json)
model.save_weights(model_dir + "resnet50_imagenet_weights.h5")
print("Saved model to disk")