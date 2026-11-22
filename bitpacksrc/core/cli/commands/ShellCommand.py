import os
import json


class ShellCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack shell <module>"
            )
            return

        modules_dir = (
            "data/modules"
        )

        target = args[0]

        if not os.path.exists(
            modules_dir
        ):
            print(
                "No modules directory."
            )
            return

        for file in os.listdir(modules_dir):

            if not file.endswith(".json"):
                continue

            path = os.path.join(
                modules_dir,
                file
            )

            with open(path) as f:
                module = json.load(f)

            if module["name"] == target:

                print(
                    f"Attached to module "
                    f"{target}"
                )

                print(
                    "Shell support "
                    "coming soon."
                )

                return

        print("Module not found.")