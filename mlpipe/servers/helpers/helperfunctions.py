import sys
import base64
import numpy as np
import json


def base64_encoding(array):
    return base64.b64encode(array).decode("utf-8")


def base64_decoding(array, dtype, shape):
    if sys.version_info.major == 3:
        array = bytes(array, encoding="utf-8")
    array = np.frombuffer(base64.decodestring(array), dtype=dtype) # decodesting() depricated to decodebytes()
    array = array.reshape(shape)
    return array


def get_dtype(array):
    return str(array.dtype)


class NumpyEncoder(json.JSONEncoder):
    """ JSON encoder for numpy types """
    def default(self, obj):
        if isinstance(obj, (
                np.int_, np.intc, np.intp, np.int8, np.int16, np.int32, np.int64, np.uint8, np.uint16, np.uint32, np.uint64)):
            return int(obj)
        elif isinstance(obj, (np.float_, np.float16, np.float32, np.float64)):
            return float(obj)
        elif isinstance(obj, (np.ndarray, )):
            return obj.tolist()
        return json.JSONEncoder.default(self, obj)