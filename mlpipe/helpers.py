import sys
import base64
import numpy as np


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