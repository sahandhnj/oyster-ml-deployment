"""
Resources
"""
import os
import yaml
from flask import Flask, flash, request, redirect, url_for, jsonify
from werkzeug.utils import secure_filename
with open("./config/settings.yaml", 'r') as stream:
    try:
        settings = yaml.load(stream)
    except yaml.YAMLError as exc:
        print(exc)


ROOT_FOLDER = os.path.abspath(".")
DATA_FOLDER = str(ROOT_FOLDER + "/data")

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = DATA_FOLDER


def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in set(settings['data_stream']['allowed_extensions'])



# curl -i -X PUT -F name=Test -F filedata=@SomeFile.pdf "http://localhost:5000/"
# curl -X POST -F "data=@review.csv" http://localhost:5000/

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


if __name__ == "__main__":
    app.run()



# dog = allowed_file('dog.jpg')
# dog2 = allowed_file('dog.xml')
#
# print(dog, dog2)