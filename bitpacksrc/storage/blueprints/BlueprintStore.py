import json


class BlueprintStore:

    ROOT = "./data/blueprints"

    def __init__(self):
        self.ROOT.mkdir(
            parents=True,
            exist_ok=True
        )

    def save(
        self,
        name,
        manifest
    ):

        path = (
            self.ROOT /
            f"{name}.bp"
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

    def load(
        self,
        name
    ):

        path = (
            self.ROOT /
            f"{name}.bp"
        )

        if not path.exists():
            return None

        with open(
            path,
            "r",
            encoding="utf-8"
        ) as file:

            return json.load(file)

    def delete(
        self,
        name
    ):
        (
            self.ROOT /
            f"{name}.bp"
        ).unlink(
            missing_ok=True
        )