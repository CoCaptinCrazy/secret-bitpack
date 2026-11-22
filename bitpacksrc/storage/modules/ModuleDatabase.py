import json


class ModuleDatabase:

    FILE =  "./data/modules/database.json"

    def __init__(self):

        self.FILE.parent.mkdir(
            parents=True,
            exist_ok=True
        )

        if not self.FILE.exists():
            self._save([])

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
        module_id
    ):

        data = self._load()

        if module_id not in data:
            data.append(
                module_id
            )

        self._save(data)

    def list(self):
        return self._load()