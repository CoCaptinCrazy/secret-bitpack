import json


class ModuleStore:

    ROOT = "./data/modules"

    def __init__(self):

        self.ROOT.mkdir(
            parents=True,
            exist_ok=True
        )

    def save(
        self,
        module
    ):

        path = (
            self.ROOT /
            f"{module.id}.json"
        )

        with open(
            path,
            "w",
            encoding="utf-8"
        ) as file:

            json.dump(
                module.to_dict(),
                file,
                indent=4
            )

    def load(
        self,
        module_id
    ):

        path = (
            self.ROOT /
            f"{module_id}.json"
        )

        if not path.exists():
            return None

        with open(
            path,
            "r",
            encoding="utf-8"
        ) as file:

            return json.load(file)