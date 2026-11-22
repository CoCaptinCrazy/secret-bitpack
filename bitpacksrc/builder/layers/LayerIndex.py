import json


class LayerIndex:

    INDEX_FILE = "./data/layers/index.json"

    def __init__(self):

        self.INDEX_FILE.parent.mkdir(
            parents=True,
            exist_ok=True
        )

        if not self.INDEX_FILE.exists():
            self._save({})

    def _load(self):

        with open(
            self.INDEX_FILE,
            "r"
        ) as file:

            return json.load(file)

    def _save(self, data):

        with open(
            self.INDEX_FILE,
            "w"
        ) as file:

            json.dump(
                data,
                file,
                indent=4
            )

    def register(
        self,
        blueprint,
        layer_id
    ):

        data = self._load()

        data.setdefault(
            blueprint,
            []
        ).append(layer_id)

        self._save(data)

    def get_layers(
        self,
        blueprint
    ):
        data = self._load()

        return data.get(
            blueprint,
            []
        )