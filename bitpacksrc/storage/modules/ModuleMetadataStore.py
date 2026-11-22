import json


class ModuleMetadataStore:

    ROOT = "./data/modules/metadata"

    def __init__(self):

        self.ROOT.mkdir(
            parents=True,
            exist_ok=True
        )

    def save(
        self,
        module_id,
        metadata
    ):

        path = (
            self.ROOT /
            f"{module_id}.json"
        )

        with open(
            path,
            "w"
        ) as file:

            json.dump(
                metadata,
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
            return {}

        with open(
            path,
            "r"
        ) as file:

            return json.load(file)