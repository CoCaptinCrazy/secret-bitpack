class MountManager:

    def __init__(self):
        self.mounts = {}

    def mount(
        self,
        module_id,
        volume_name,
        target
    ):

        self.mounts.setdefault(
            module_id,
            []
        ).append({
            "volume": volume_name,
            "target": target
        })

    def get(
        self,
        module_id
    ):
        return self.mounts.get(
            module_id,
            []
        )