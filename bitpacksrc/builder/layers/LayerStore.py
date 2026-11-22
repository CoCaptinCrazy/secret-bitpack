import json

class LayerStore:

    ROOT = "./data/layers"

    def __init__(self):
        self.ROOT.mkdir(
            parents=True,
            exist_ok=True
        )

    def save(self, layer):

        path = (
            self.ROOT /
            f"{layer.id}.json"
        )

        with open(
            path,
            "w",
            encoding="utf-8"
        ) as file:

            json.dump(
                layer.to_dict(),
                file,
                indent=4
            )

    def exists(self, layer_id):

        return (
            self.ROOT /
            f"{layer_id}.json"
        ).exists()