"""
Resources
"""

import os
from flask import Flask, flash, request, redirect, url_for, jsonify
from werkzeug.utils import secure_filename

ROOT_FOLDER = os.path.abspath(".")
DATA_FOLDER = str(ROOT_FOLDER + "/data")
ALLOWED_EXTENSIONS = set(['txt', 'png', 'jpg', 'jpeg', 'json'])

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = DATA_FOLDER


def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS



# curl -i -X PUT -F name=Test -F filedata=@SomeFile.pdf "http://localhost:5000/"
# curl -X POST -F "data=@tweet.txt" http://localhost:5000/

@app.route("/", methods=["POST"])
def hello():
    if request.method == "POST":
        # Check if file is inputted
        if 'data' not in request.files:
            flash('No file part')
            raise ValueError("No file part 2")
        file = request.files['data']
        # Check if filename is not empty
        if file.filename == '':
            flash('No selected file')
        if file and allowed_file(file.filename):
            filename=secure_filename(file.filename)
            print(filename)
        if request.files.get('data'):
            user_input = request.files["data"].read()
            print(user_input)
            file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))



    return "Success"




# @app.route('/perdict', methods=["GET", "POST"])
# def upload_file():
#     if request.method == "POST":
#         if request.files.get("Test"):
#             file = request.files["Test"].read()
#             if file and allowed_file(file.filename):
#                 filename = secure_filename(file.filename)
#                 print(filename)
#
#             return "Success"



    # return jsonify(user_input)

            #
    # if 'file' not in request.files:
    #     flash('No file part')
    #     raise ValueError("No file part")
    # file = request.files['file']
    # if file.filename == '':
    #     flash('No selected file')
    #     raise ValueError("No selected file")
    # if file and allowed_file(file.filename):
    #     filename = secure_filename(file.filename)
    #     user_input = request.files['data'].read()
    #
    #     # file.save(os.path.join(app.config['UPLOAD_FOLDER'], filename))
    #     print("File uploaded to folder", filename, user_input)
    # print("File uploaded to folder", filename, user_input)





if __name__ == "__main__":
    app.run()



# dog = allowed_file('dog.jpg')
# dog2 = allowed_file('dog.xml')
#
# print(dog, dog2)