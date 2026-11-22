import json


class VolumeStore:

    FILE = "./data/volumes.json"

    def load(self):

        if not self.FILE.exists():
            return {}

        with open(
            self.FILE,
            "r"
        ) as file:

            return json.load(file)

    def save(
        self,
        data
    ):

        with open(
            self.FILE,
            "w"
        ) as file:

            json.dump(
                data,
                file,
                indent=4
            )