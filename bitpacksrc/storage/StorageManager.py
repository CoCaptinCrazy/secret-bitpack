import os


class StorageManager:

    ROOT = "./data"

    PATHS = {
        "blueprints": ROOT + "/blueprints",
        "modules": ROOT + "/modules",
        "layers": ROOT + "/layers",
        "volumes": ROOT + "/volumes",
        "caches": ROOT + "/caches",
        "meshes": ROOT + "/meshes",
        "assemblies": ROOT + "/assemblies",
        "clusters": ROOT + "/clusters",
        "logs": ROOT + "/logs",
        "snapshots": ROOT + "/snapshots",
        "backups": ROOT + "/backups"
    }

    def initialize(self):

        os.makedirs(
            self.ROOT,
            exist_ok=True
        )

        for path in self.PATHS.values():

            os.makedirs(path, exist_ok=True)

    def path(self, name):
        return self.PATHS[name]