import csv
import os

with open(os.path.abspath("./tests") + "/review.csv", "r") as f:
    reader = csv.reader(f, delimiter=',')

    tmp = []
    for i, line in enumerate(reader):
        sentence = line[0]
        print(sentence)

        for word in sentence.split(" "):
            tmp.append(word)
    print(tmp)