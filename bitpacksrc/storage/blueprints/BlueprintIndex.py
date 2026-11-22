import json


class BlueprintIndex:

    FILE = "./data/blueprints/index.json"

    def __init__(self):

        self.FILE.parent.mkdir(
            parents=True,
            exist_ok=True
        )

        if not self.FILE.exists():
            self._save({})
            
    def _load(self):

        with open(
            self.FILE,
            "r"
        ) as file:
            return json.load(file)

    def _save(
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

    def register(
        self,
        name,
        version
    ):

        data = self._load()

        data[name] = {
            "version": version
        }

        self._save(data)

    def list(self):
        return self._load()