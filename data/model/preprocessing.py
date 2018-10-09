"""
Users can upload their own preprocessing file
Requirement:
- File has to contain a function called preprocessing
"""


import numpy
from keras.preprocessing import sequence
import keras

numpy.random.seed(42)

max_review_length = 500
top_words = 5000
NUM_WORDS = 1000  # only use top 1000 words
INDEX_FROM = 3  # word index offset

word_to_id = keras.datasets.imdb.get_word_index()
word_to_id = {k: (v + INDEX_FROM) for k, v in word_to_id.items()}
word_to_id["<PAD>"] = 0
word_to_id["<START>"] = 1
word_to_id["<UNK>"] = 2

id_to_word = {value: key for key, value in word_to_id.items()}

def preprocessing(data):
    tmp = []
    for word in data.split(" "):
        tmp.append(word_to_id[word])
    tmp_padded = sequence.pad_sequences([tmp], maxlen=max_review_length)
    return tmp_padded