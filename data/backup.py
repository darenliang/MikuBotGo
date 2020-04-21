import requests

if __name__ == "__main__":
    folders = list(range(2000, 2021)) + ["60s", "70s", "80s", "90s", "misc"]

    for i in folders:
        r = requests.get(f"https://aninx.com/{i}/success.txt")
        with open(f"{i}success.txt", 'wb') as f:
            f.write(r.content)
