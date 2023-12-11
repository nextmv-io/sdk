import json
import math

with open("input.json", "r", encoding="utf-8") as file:
    input = json.load(file)

for demand in input["demands"]:
    demand["demand"] = int(math.ceil(demand["demand"] / 10.0))

with open("input-reduced.json", "w", encoding="utf-8") as file:
    json.dump(input, file, indent=2)
