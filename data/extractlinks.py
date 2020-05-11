import json

with open('dataset_filtered_old.json') as f:
    data = json.load(f)

for anime in data:
    for song in anime['songs']:
        print(song['url'])
