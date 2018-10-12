"""
Command line interface to store user input into
"""

import os
import sys
import h5py
import json
from datetime import date, datetime


cli_params = {
    "MIN_ID_CHAR": 5,
    "MAX_ID_CHAR": 32
}


def is_json(json_file_path):
    try:
        with open(json_file_path, 'r') as jsonfile:
            _model = jsonfile.read()
        return True
    except ValueError:
        print("Not a valid JSON file.")
        return False


def json_serial(obj):
    """JSON serializer for objects not serializable by default json code"""

    if isinstance(obj, (datetime, date)):
        return obj.isoformat()
    raise TypeError("Type %s not serializable" % type(obj))


def is_hdf5(file_name):
    try:
        # TBD: Wrap in with-open clause
        _data = h5py.File("{}".format(file_name), 'r')
        return True
    except ValueError:
        print("Not a valid HDF5 file.")
        return False


def get_user_input():

    print("Welcome to Oysterbox!\n"
          "Specify the pearls you like to ship over the web.\n")

    # Question 1: Project name
    while True:
        project_id = input("1. Project name?\n")
        if len(project_id) < cli_params["MIN_ID_CHAR"]:
            print("Name is to short, minimum is 5 characters.\nPlease try again.")
        elif len(project_id) > cli_params["MAX_ID_CHAR"]:
            print("Name is to long, maximum is 32 characters.\nPlease try again.")
        else:
            break

    # Question 2: Base module / main dependency
    while True:
        module = input("2. Base module? (Currently only supporting 'keras')\n")
        if module == 'keras':
            break
        else:
            print("Please choose: keras | pytorch | sklearn.")

    # Question 3: Computation graph file
    while True:
        graph = input("3. Computation graph file (JSON)?\n")
        graph_file_exists = os.path.isfile(graph)
        if graph_file_exists:
            json_exists = is_json(graph)
            if json_exists:
                break
            else:
                print("Not a valid JSON file. Please select again.")
        else:
            print("That's not an existing file. Please check your path.")

    # Question 4: Model weights file
    while True:
        weights = input("4. Model weights file (HDF5)?\n")
        weights_file_exists = os.path.isfile(weights)
        if weights_file_exists:
            hdf5_exists = is_hdf5(weights)
            if hdf5_exists:
                break
            else:
                print("Not a valid HDF5 file. Pleas select again.")
        else:
            print("This is not an existing file. Please check your path.")

    # Question 5: Description (optional)
    description = input("5. Description (optional)\n")

    # Default naming
    project_id = 'project1'
    module = 'keras',
    graph = 'graph.json',
    weights = 'weights.h5',
    description = 'description'

    return project_id, module, graph, weights, description


def create_metadata_file(project_id, module, graph, weights, description):

    current_time = json_serial(datetime.now())
    metadata = {
        "project_id": project_id,
        "module": module,
        "created": current_time,
        "graph_file": graph,
        "weights_file": weights,
        "description": description
    }

    with open("meta.json", 'w') as metafile:
        json.dump(metadata, metafile)

    with open("meta.json", 'r') as fp:
        data = json.load(fp)

    print("Meta data written to /model/meta.json, with following content:\n")
    print(data)


def main():

    file_id, module, graph, weights, description = get_user_input()
    create_metadata_file(file_id, module, graph, weights, description)
    print("Thank you, creating vessel for shipment...")

    return file_id, module, graph, weights, description


if __name__ == "__main__":

    main()
