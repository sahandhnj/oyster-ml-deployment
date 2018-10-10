import time
import uuid
import json
import inspect
import yaml
import redis
from config.clistyle import bcolor
from flask import Flask, request, jsonify
from helpers import base64_encoding, get_dtype

try:
    from model import preprocessing as prepmod

    if hasattr(prepmod, 'preprocessing') and inspect.isfunction(prepmod.preprocessing):
        from model.preprocessing import preprocessing
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

# Flask variables
# ALLOWED_EXTENSIONS = set(['txt', 'png', 'jpg', 'jpeg', 'wav'])

# rdb.flushall()


@app.route('/predict', methods=["POST"])
def predict():

    data = {"success": False}

    if request.method == "POST":
        user_input = request.json["text"]      
        preprocessed_input = preprocessing(user_input)
        array_dtype = get_dtype(preprocessed_input)            
        encoded_input = base64_encoding(preprocessed_input)
        # endoced_input = encoded_input.copy(order="C")   # make C-contigious?
        
        k = str(uuid.uuid4())
        d = {"id": k, "shape": preprocessed_input.shape, "dtype": array_dtype, "data": encoded_input}
        rdb.rpush(settings['data_stream']['data_queue'], json.dumps(d))    # dump the preprocessed input as a numpy array

        while True:
            output = rdb.get(k)

            if output is not None:
                output = output.decode("utf-8")
                data["predictions"] = json.loads(output)

                rdb.delete(k)
                break
            
            time.sleep(settings['data_stream']['client_sleep'])
        data["success"] = True
    
    return jsonify(data)    


@app.route("/")
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
    app.run(
        host=settings['webservice']['host'],
        port=int(settings['webservice']['port']),
        debug=settings['webservice']['debug']
    )