import os, sys, time, json, yaml, uuid
import redis
import torch
import numpy as np
from PIL import Image
from werkzeug.utils import secure_filename
from flask import Flask, request, flash, jsonify
mlpipe_root = os.path.abspath("../..")
# sys.path.insert(1, os.path.join(sys.path[0], mlpipe_root)
sys.path.insert(0, mlpipe_root)

# for p in sys.path:
#     print(p + "\n")
# print(mlpipe_root)

from config.clistyle import bcolor
from servers.helpers.helperfunctions import base64_encoding, get_dtype

# Preprocessing specific
from torchvision import transforms

with open(mlpipe_root + "/config/settings.yaml", 'r') as stream:
    try:
        settings = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)

# TBD after merge
# with open("./config/allowedExtns.yaml", 'r') as stream:
#     try:
#         allowed_extensions = yaml.load(stream)
#     except yaml.YAMLError as exc:
#         print(exc)

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


# def get_device():
#     return torch.device("cuda:0" if torch.cuda.is_available() else "cpu")


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


@app.route("/predict", methods=["POST"])
def predict():
    
    # device = get_device()
    data = {"success": False}

    if request.method == "POST":
        # Check if file in inputted
        if 'data' not in request.files:
            flash("No file part")
            raise ValueError("No file part")
        # print("METHOD: ", request.method)
        file = request.files['data']
        filetype = get_file_type(file.filename)
        if file.filename == '':
            flash("No selected file")
            raise ValueError("No selected file")
        if file and allowed_file(file.filename):
            filename = secure_filename(file.filename)
            if request.files.get('data'):
                user_input = request.files["data"]  #.read()
                if (filetype in ['jpg', 'jpeg', 'png']):
                    user_input = Image.open(user_input)
                else:
                    pass
        
                preprocessed_input = preprocessing(user_input)               
                # Get file properties
                if filetype in ['jpg', 'jpeg', 'png']:
                    fileshape = np.array(preprocessed_input).shape
                else:
                    fileshape = preprocessed_input.shape

                array_dtype = get_dtype(preprocessed_input)
                # HANDLE TORCH TENSORS
                # print("DTYPE: ", array_dtype)
                # if array_dtype == 'torch._':
                #     print("TORCH!")
                #     ### convert torch to numpy
                #     output_np = torchTensor_to_npArray(output1)
                
                preprocessed_input = preprocessed_input.numpy()         
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
                rdb.rpush(settings['data_stream']['data_queue'], json.dumps(d))

                # Can i also send and receive torch tensors via redis?
                # or setup automated function to resore them from np array
                while True:
                    output = rdb.get(k)
                    if output is not None:
                        output = output.decode("utf-8")
                        data["summary"] = json.loads(output)
                        rdb.delete(k)
                        break
                    
                    time.sleep(settings['data_stream']['client_sleep'])
                data["success"] = True

    return jsonify(data)    
       

if __name__ == "__main__":
    print((bcolor.BOLD + "* Loading PyTorch webserver... \n"
           "please wait until server has fully started" + bcolor.END))   
    print("* Starting web service...")
    app.secret_key = settings['flask']['secret_key']
    app.run(
        host=settings['flask']['host'],
        port=int(settings['flask']['port']),
        debug=settings['flask']['debug']
    )