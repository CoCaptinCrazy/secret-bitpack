import json
import os


class RestartCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack restart <module>"
            )
            return

        target = args[0]

        modules = (
            "data/modules"
        )

        for file in os.listdir(modules):

            path = os.path.join(
                modules,
                file
            )

            with open(path) as f:
                module = json.load(f)

            if (
                module["name"] == target
            ):

                module["status"] = (
                    "running"
                )

                with open(
                    path,
                    "w"
                ) as f:

                    json.dump(
                        module,
                        f,
                        indent=4
                    )

                print(
                    f"Restarted "
                    f"{target}"
                )

                return

        print(
            "Module not found."
        )