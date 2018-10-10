
import tensorflow as tf
from keras.models import model_from_json

from keras.preprocessing import sequence
import keras

# numpy.random.seed(42)

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


def load_model(model_file_path, weights_file_path):
    global model
    with open("{}".format(model_file_path), 'r') as model_json_file:
        loaded_model_json = model_json_file.read()
    loaded_model = model_from_json(loaded_model_json)
    loaded_model.load_weights("{}".format(weights_file_path))

    global graph
    graph = tf.get_default_graph()
    print("Loaded model from disk and inserted weights.")

    return loaded_model

model = load_model('model.json', 'model_weights.h5')

user_inp = "you hate the movie" # you love the movie
processed_inp = preprocessing(user_inp)

def get_dtype(array):
    return str(array.dtype)

d = get_dtype(processed_inp)
print(d)

print(processed_inp)


with graph.as_default():
    predictions = model.predict(processed_inp)
    print(predictions)

# print(predictions)