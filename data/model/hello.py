import tensorflow as tf
import numpy
from keras.preprocessing import sequence
from keras.models import model_from_json
import keras
from flask import Flask, request

numpy.random.seed(7)

app = Flask(__name__)


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



# Data preprocessing
def word2id(word):
    INDEX_FROM = 3  # word index offset
    word2id_list = keras.datasets.imdb.get_word_index()
    word2id_list = {k: (v + INDEX_FROM) for k, v in word2id_list.items()}
    word2id_list["<PAD>"] = 0
    word2id_list["<START>"] = 1
    word2id_list["<UNK>"] = 2
    word_id = word2id_list[word]

    return word_id


def sequence_padding(sentence, maxlen=500):
    tmp = []
    for word in sentence.split(" "):
        tmp.append(word2id(word))
    tmp_padded = sequence.pad_sequences([tmp], maxlen=maxlen)

    return tmp_padded


def preprocessing(sentence):
    prepped_data = sequence_padding(sentence)
    return prepped_data

# Perform all steps
# 1. Load model
model = load_model('model.json', 'model_weights.h5')

@app.route('/predict', methods=['POST'])
def predict():
    # 2. retrieve input
    user_input = request.json['text']
    #user_input = "n the movie was terrible"
    # print(user_input)
    # 3. preprocess data
    data = preprocessing(user_input[0])
    # print(data)

    # 4. output results (should become api)
    with graph.as_default():
        preds = model.predict(data)

    print("Sentiment: {}".format(preds))
    return "Sentiment: {}".format(preds)


@app.route("/")
def hello():
    return "Hello World!"

if __name__ == "__main__":
    print(("* Loading Keras model and Flask starting server..."
        "please wait until server has fully started"))
    app.run(host="0.0.0.0", port=int("5000"), debug=True)