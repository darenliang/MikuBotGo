import os
from os import listdir
from os.path import isfile, join

import youtube_dl

existing_files = [f for f in listdir("D:/TMP") if isfile(join("D:/TMP", f))]

print(f"Number of existing files: {len(existing_files)}")

with open('links.txt', 'r') as f:
    entries = f.readlines()

for i in range(len(entries)):
    entries[i] = entries[i].rstrip()

print(f"Number of entries in database: {len(entries)}")

ydl_opts = {
    'format': 'bestaudio/best',
    'outtmpl': "D:/TMP/%(title)s.%(ext)s",
    'postprocessors': [{
        'key': 'FFmpegExtractAudio',
        'preferredcodec': 'mp3',
        'preferredquality': '192',
    }],
}

for existing_file in existing_files:
    if existing_file not in entries:
        input(f"Enter to remove {existing_file}")
        print(f"Removing {existing_file}")
        os.remove(f"D:/TMP/{existing_file}")

for entry in entries:
    if entry not in existing_files:
        with youtube_dl.YoutubeDL(ydl_opts) as ydl:
            ydl.download([f"https://animethemes.moe/video/{entry[:-3]}webm"])

print("Done.")
