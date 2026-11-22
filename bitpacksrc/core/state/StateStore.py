import json
import os


class StateStore:

    FILE = "./data/state.json"

    def __init__(self):

        os.makedirs(
            os.path.dirname(self.FILE),
            exist_ok=True
        )

    def load(self):

        if not os.path.exists(self.FILE):
            return {}

        with open(
            self.FILE,
            "r",
            encoding="utf-8"
        ) as file:

            return json.load(file)

    def save(self, data):

        with open(
            self.FILE,
            "w",
            encoding="utf-8"
        ) as file:

            json.dump(
                data,
                file,
                indent=4
            )