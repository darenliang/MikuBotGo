import json
import re
import os


class Aninx:
    Endpoint = "https://aninx.com"
    Folder = ""
    Data = []


def remove_prefix(text, prefix):
    return text[text.startswith(prefix) and len(prefix):].rstrip("\n")


if __name__ == "__main__":
    # for i in range(2000, 2021):
    #     r = requests.get(f"{Aninx.Endpoint}/{i}/success.txt")
    #     with open(f"{Aninx.Folder}/{i}success.txt", 'wb') as f:
    #         f.write(r.content)

    for i in range(2000, 2021):
        Aninx.Data.append({"year": i, "animes": []})
        with open(f"{Aninx.Folder}{i}success.txt", "r", encoding="utf8", errors="ignore") as f:
            line = f.readline()
            while line:
                search = re.search("^.{7}─ (.+)\n$", line)
                if search:
                    Aninx.Data[-1]["animes"].append({"name": [search[1]], "songs": []})
                else:
                    search = re.search("^.{12}└─ 0: (.+)\n$", line)
                    if search:
                        Aninx.Data[-1]["animes"][-1]["songs"].append(search[1])
                line = f.readline()

    with open(f"{Aninx.Folder}dataset.json", "w") as f:
        json.dump(Aninx.Data, f)
