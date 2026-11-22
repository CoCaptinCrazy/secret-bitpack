import os
import json
import time
import shutil


class SnapshotCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bitpack snapshot <module>"
            )

            return

        module = args[0]

        module_path = os.path.join(
            "data",
            "modules",
            module
        )

        if not os.path.exists(module_path):

            print(
                f"Module not found: {module}"
            )

            return

        snapshot_id = str(
            int(time.time())
        )

        snapshot_path = os.path.join(
            "data",
            "snapshots",
            module,
            snapshot_id
        )

        os.makedirs(
            os.path.dirname(snapshot_path),
            exist_ok=True
        )

        shutil.copytree(
            module_path,
            snapshot_path
        )

        metadata = {
            "module": module,
            "snapshot": snapshot_id,
            "created": int(time.time())
        }

        with open(
            os.path.join(
                snapshot_path,
                "snapshot.json"
            ),
            "w"
        ) as file:

            json.dump(
                metadata,
                file,
                indent=4
            )

        print(
            f"Snapshot created: {snapshot_id}"
        )