import json
import os
import datetime


class InspectModuleCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bitpack inspect-module <module>"
            )

            return

        target = args[0]

        modules_dir = (
            "data/modules"
        )

        if not os.path.exists(
            modules_dir
        ):

            print(
                "No modules found."
            )

            return

        found = None

        for file in os.listdir(
            modules_dir
        ):

            if not file.endswith(
                ".json"
            ):
                continue

            path = os.path.join(
                modules_dir,
                file
            )

            with open(
                path,
                "r",
                encoding="utf-8"
            ) as module_file:

                module = json.load(
                    module_file
                )

            if (
                module.get("name")
                == target
                or
                module.get("id")
                == target
            ):

                found = module

                break

        if not found:

            print(
                f"Module not found: {target}"
            )

            return

        started = found.get(
            "started"
        )

        uptime = "Unknown"

        if started:

            now = (
                datetime.datetime.now()
                .timestamp()
            )

            uptime = (
                f"{int(now - started)}s"
            )

        print()

        print(
            "MODULE INSPECTION"
        )

        print(
            "=" * 40
        )

        print(
            f"ID:        {found.get('id')}"
        )

        print(
            f"NAME:      {found.get('name')}"
        )

        print(
            f"STATUS:    {found.get('status')}"
        )

        print(
            f"BLUEPRINT: "
            f"{found.get('blueprint')}"
        )

        print(
            f"RUNTIME:   "
            f"{found.get('runtime')}"
        )

        print(
            f"UPTIME:    "
            f"{uptime}"
        )

        print()

        snapshots_dir = (
            "data/snapshots"
        )

        snapshot_count = 0

        if os.path.exists(
            snapshots_dir
        ):

            snapshot_count = len(
                os.listdir(
                    snapshots_dir
                )
            )

        print(
            f"SNAPSHOTS: "
            f"{snapshot_count}"
        )

        print()