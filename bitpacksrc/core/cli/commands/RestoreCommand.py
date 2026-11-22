import os
import shutil


class RestoreCommand:

    def execute(self, args):

        if len(args) < 2:

            print(
                "Usage: bitpack restore <module> <snapshot>"
            )

            return

        module = args[0]

        snapshot = args[1]

        snapshot_path = os.path.join(
            "data",
            "snapshots",
            module,
            snapshot
        )

        if not os.path.exists(snapshot_path):

            print(
                f"Snapshot not found: {snapshot}"
            )

            return

        module_path = os.path.join(
            "data",
            "modules",
            module
        )

        if os.path.exists(module_path):

            shutil.rmtree(
                module_path
            )

        shutil.copytree(
            snapshot_path,
            module_path
        )

        snapshot_file = os.path.join(
            module_path,
            "snapshot.json"
        )

        if os.path.exists(snapshot_file):

            os.remove(snapshot_file)

        print(
            f"Module restored: {module}"
        )