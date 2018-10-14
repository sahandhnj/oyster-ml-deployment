import time
import uuid
import json
import inspect
import yaml
import redis
from config.clistyle import bcolor
from flask import Flask, request, jsonify, flash
from werkzeug.utils import secure_filename
from helpers import base64_encoding, get_dtype
from PIL import Image
import io
import numpy as np

try:
    from runs.imagenet import preprocessing as prepmod

    if hasattr(prepmod, 'preprocessing') and inspect.isfunction(prepmod.preprocessing):
        from runs.imagenet.preprocessing import preprocessing
        print("Preprocessing file available and loaded into vessel.")
    else:
        raise TypeError("Preprocessing file inserted, but does not contain function called 'preprocessing'.")
except (ImportError):
    print("No preprocessing file inserted.")

with open("./config/settings.yaml", 'r') as stream:
    try:
        settings = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)



# numpy.random.seed(42)
app = Flask(__name__)
rdb = redis.StrictRedis(
    host=settings['redis']['host'],
    port=settings['redis']['port'],
    db=settings['redis']['db']
)

# rdb.flushall()


def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in set(settings['data_stream']['allowed_extensions'])


def get_file_type(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower()


@app.route('/predict', methods=["POST"])
def predict():

    data = {"success": False}

    if request.method == "POST":
        # Check if file is inputted
        if 'data' not in request.files:
            flash("No file part")
            raise ValueError("No file part.")
        file = request.files['data']
        # print("FILENAME: ", file.filename)
        filetype = get_file_type(file.filename)
        # print("FILETYPE: ", filetype)
        # Check if file name is not empty
        if file.filename == '':
            flash("No selected file")
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
        if request.files.get('data'):
            user_input = request.files["data"].read()
            if (filetype in ['jpg', 'jpeg', 'png']): 
                user_input = Image.open(io.BytesIO(user_input))
            else:
                pass

            preprocessed_input = preprocessing(user_input)           
            # Get file properties
            if filetype in ['jpg', 'jpeg', 'png']: 
                fileshape = np.array(preprocessed_input).shape
            else:
                fileshape = preprocessed_input.shape
      
            array_dtype = get_dtype(preprocessed_input)
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
            rdb.rpush(settings['data_stream']['data_queue'], json.dumps(d))  # dump the preprocessed input as a numpy array

            while True:
                output = rdb.get(k)
                if output is not None:
                    output = output.decode("utf-8")
                    # print("SUMMARY: ", json.loads(output)[0])
                    data["summary"] = json.loads(output)
                    rdb.delete(k)
                    break

                time.sleep(settings['data_stream']['client_sleep'])
            data["success"] = True
   
    return jsonify(data)    


@app.route("/predict")
def hello():
    return "Hello, Welcome to Oysterbox Machine Learning Deployment!"


if __name__ == "__main__":
    print((bcolor.BOLD + "* Loading Keras model and Flask starting server... \n"
           "please wait until server has fully started" + bcolor.END))
    # print("* Starting model service... ")
    # t = Thread(target=classify_process, args=())
    # t.daemon
    # t.start()
    
    print("* Starting web service...")
    app.secret_key = settings['flask']['secret_key']
    # app.config['SESSION_TYPE'] = 'filesystem'
    app.run(
        host=settings['flask']['host'],
        port=int(settings['flask']['port']),
        debug=settings['flask']['debug']
    )