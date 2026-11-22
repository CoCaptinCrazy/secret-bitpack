class VolumeDriver:

    ROOT = "./data/volumes"

    def __init__(self):
        self.ROOT.mkdir(
            parents=True,
            exist_ok=True
        )

    def create(
        self,
        volume
    ):

        path = (
            self.ROOT /
            volume.name
        )

        path.mkdir(
            parents=True,
            exist_ok=True
        )

        return path

    def remove(
        self,
        volume_name
    ):
        path = (
            self.ROOT /
            volume_name
        )

        if path.exists():
            path.rmdir()