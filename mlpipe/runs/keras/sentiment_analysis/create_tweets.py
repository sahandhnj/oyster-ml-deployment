
import os
import uuid
import csv

cwd = os.getcwd()
data_dir = cwd + "/tweets/"

tweet_files = os.listdir(data_dir)
tweet_id = str(uuid.uuid4())[:6]
filename = "tweet"

tweet_list = [
    'you like the movie',
    'you hate the movie',
    'they love their job',
    'we like the movie',
    'they play nice',
    'we did not understand',
    'we are cool',
    'movie was awesome',
    'movie was terrible',
    'we think it was cool',
    'we are awesome',
    'they are awesome',
    'there are terrible',
    'good job',
    'nice job',
    'terrible job',
    'cool job',
    'great job',
    'bad movie',
    'we love bad movies',
    'we love nice movies',
    'we hate bad movies',
    'we hate nice movies',
    'they know bad movies',
    'they do nice job',
    'they do bad job',
    'they understand it perfect',
    'they are bad',
    'they are nice',
    'they are awesome',
    'they like bad movies',
    'bad movies are a good job',
    'movie was well',
    'movies are bad',
    'movies are dangerous',
    'very nice movie',
    'very nice movies',    
]

for tweet in tweet_list:
    with open(data_dir + "review" + str(uuid.uuid4())[:6] + ".csv", "w") as csvfile:
        csvfile.write(tweet)
                # writer = csv.writer(csvfile)
        # writer.writerows(tweet)


