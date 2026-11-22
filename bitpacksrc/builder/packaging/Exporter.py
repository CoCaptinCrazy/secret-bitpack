import json


class Exporter:

    OUTPUT = "./data/blueprints"

    def __init__(self):
        self.OUTPUT.mkdir(
            parents=True,
            exist_ok=True
        )

    def export(
        self,
        blueprint,
        manifest
    ):

        path = (
            self.OUTPUT /
            f"{blueprint.name}.bp"
        )

        with open(
            path,
            "w",
            encoding="utf-8"
        ) as file:

            json.dump(
                manifest,
                file,
                indent=4
            )

        return path