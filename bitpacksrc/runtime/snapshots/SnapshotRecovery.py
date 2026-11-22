class SnapshotRecovery:

    def recover(
        self,
        snapshot_data
    ):

        return {
            "module_id":
                snapshot_data["module_id"],

            "state":
                snapshot_data["state"]
        }