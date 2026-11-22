import json


class CacheDriver:

    ROOT = "./data/caches"

    def __init__(self):
        self.ROOT.mkdir(
            parents=True,
            exist_ok=True
        )

    def save(
        self,
        cache
    ):
        path = (
            self.ROOT /
            f"{cache.name}.json"
        )

        with open(
            path,
            "w",
            encoding="utf-8"
        ) as file:

            json.dump(
                cache.data,
                file,
                indent=4
            )

    def load(
        self,
        name
    ):
        path = (
            self.ROOT /
            f"{name}.json"
        )

        if not path.exists():
            return {}

        with open(
            path,
            "r",
            encoding="utf-8"
        ) as file:

            return json.load(file)