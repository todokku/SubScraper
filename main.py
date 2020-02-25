#!/usr/bin/python3

import requests,json
import sys
word = sys.argv[1]

resp = requests.get("https://api.dictionaryapi.dev/api/v1/entries/en/"+word)

data = json.loads(resp.content)
#print(type(data[0]))

for k,v in data[0]["meaning"].items():

    output = v[0]["definition"]
    break

print(output)


#["noun"][0]["definition"])


