import json


class SnapshotStore:

    ROOT = "./data/snapshots"

    def __init__(self):
        self.ROOT.mkdir(
            parents=True,
            exist_ok=True
        )

    def save(self, snapshot):

        path = (
            self.ROOT /
            f"{snapshot.id}.json"
        )

        with open(
            path,
            "w",
            encoding="utf-8"
        ) as file:

            json.dump(
                snapshot.to_dict(),
                file,
                indent=4
            )

        return path

    def load(self, snapshot_id):

        path = (
            self.ROOT /
            f"{snapshot_id}.json"
        )

        if not path.exists():
            return None

        with open(
            path,
            "r",
            encoding="utf-8"
        ) as file:

            return json.load(file)