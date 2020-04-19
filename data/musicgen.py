import json
import pickle
import re

from fuzzywuzzy import process
from jikanpy import Jikan


class Aninx:
    Endpoint = "https://aninx.com"
    Data = []


def remove_prefix(text, prefix):
    return text[text.startswith(prefix) and len(prefix):].rstrip("\n")


def get_top_500():
    data = set()
    jikan = Jikan()

    for i in range(1, 11):
        print(f"Processing page: {i}")
        top_anime = jikan.top(type='anime', page=i, subtype='tv')
        for j in range(0, 50):
            data.add(top_anime["top"][j]["title"])
        top_anime = jikan.top(type='anime', page=i, subtype='bypopularity')
        for j in range(0, 50):
            data.add(top_anime["top"][j]["title"])

    return data


if __name__ == "__main__":
    folders = list(range(2000, 2021)) + ["60s", "70s", "80s", "90s", "misc"]
    #
    # for i in folders:
    #     r = requests.get(f"{Aninx.Endpoint}/{i}/success.txt")
    #     with open(f"{i}success.txt", 'wb') as f:
    #         f.write(r.content)

    for i in folders:
        with open(f"{i}success.txt", "r", encoding="utf8", errors="ignore") as f:
            line = f.readline()
            while line:
                search = re.search("^.{7}─ (.+)\n$", line)
                if search:
                    Aninx.Data.append({"name": search[1], "songs": []})
                else:
                    search = re.search("^.{10}─ (.+)\n$", line)
                    if search:
                        Aninx.Data[-1]["songs"].append({"songname": search[1]})
                    else:
                        search = re.search("^.{12}└─ 0: (.+)\n$", line)
                        if search:
                            Aninx.Data[-1]["songs"][-1]["url"] = search[1]
                line = f.readline()

    # filter_set = get_top_500()
    # with open('filter_set.pickle', 'wb') as handle:
    #     pickle.dump(filter_set, handle, protocol=pickle.HIGHEST_PROTOCOL)

    with open('filter_set.pickle', 'rb') as handle:
        filter_set = pickle.load(handle)

    new_data = []
    count = 0
    print(f"Total: {len(Aninx.Data)}")
    for anime in Aninx.Data:
        count += 1
        result = process.extract(anime['name'], filter_set, limit=1)[0]
        if result[1] == 100:
            new_data.append(anime)
        else:
            if result[0].endswith(" (TV)") and result[1] >= 90:
                new_data.append(anime)
        if count % 200 == 0:
            print(count)

    with open(f"dataset_filtered.json", "w") as f:
        json.dump(new_data, f)
