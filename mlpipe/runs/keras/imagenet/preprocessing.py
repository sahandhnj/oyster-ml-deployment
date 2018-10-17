"""
Users can upload their own preprocessing file
Requirement:
- File has to contain a function called preprocessing
"""


import numpy as np
from keras.preprocessing import sequence
from keras.preprocessing.image import img_to_array
from keras.applications import imagenet_utils
import os
from PIL import Image
import keras
import csv
import io
import imageio

# numpy.random.seed(42)

root = os.path.abspath(".")
print(root + "")
imagepath = root + "/tests/dog.jpg"


IMAGE_WIDTH = 224
IMAGE_HEIGHT = 224
IMAGE_CHANS = 3
IMAGE_DTYPE = "float32"

# image = Image.open(imagepath)
# image.show()
#
# image = Image.open(io.BytesIO(imagepath))   # Flask requests convert to bytesarray
# img2 = np.array(image)
# img2.shape


def prepare_image(image, target):
    if image.mode != "RGB":
        image = image.convert("RGB")
    image = image.resize(target)
    image = img_to_array(image)
    image = np.expand_dims(image, axis=0)
    image = imagenet_utils.preprocess_input(image)

    return image

# height, width, channels = imageio.imread(imagepath).shape

def preprocessing(data):

    preprocessed_data = prepare_image(data, target=(IMAGE_WIDTH, IMAGE_HEIGHT))

    return preprocessed_data




# max_review_length = 500
# top_words = 5000
# NUM_WORDS = 1000  # only use top 1000 words
# INDEX_FROM = 3  # word index offset
#
# word_to_id = keras.datasets.imdb.get_word_index()
# word_to_id = {k: (v + INDEX_FROM) for k, v in word_to_id.items()}
# word_to_id["<PAD>"] = 0
# word_to_id["<START>"] = 1
# word_to_id["<UNK>"] = 2
#
# id_to_word = {value: key for key, value in word_to_id.items()}
#
# def preprocessing(data):
#     print("DATA: ", data)
#     print("DATA: ", data.decode('utf-8'))
#     data = data.decode('utf-8')
#     tmp = []
#     # for sentence in data.split(","):
#     #     print("SENT: ", sentence)
#
#     for word in data.split(" "):
#         tmp.append(word_to_id[word])
#         tmp_padded = sequence.pad_sequences([tmp], maxlen=max_review_length)
#
#     return tmp_padded







 
    

