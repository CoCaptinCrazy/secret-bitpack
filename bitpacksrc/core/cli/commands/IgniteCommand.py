import json
import os
import time
import uuid


class IgniteCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bitpack ignite <blueprint>"
            )

            return

        blueprint_name = args[0]

        blueprint_path = os.path.join(
            "data",
            "blueprints",
            f"{blueprint_name}.bp"
        )

        if not os.path.exists(
            blueprint_path
        ):

            print(
                f"Blueprint not found: {blueprint_name}"
            )

            return

        try:

            with open(
                blueprint_path,
                "r",
                encoding="utf-8"
            ) as file:

                blueprint = json.load(file)

            metadata = blueprint.get(
                "metadata",
                {}
            )

            runtime = blueprint.get(
                "runtime",
                {}
            )

            module_id = (
                "mod-" +
                str(uuid.uuid4())[:8]
            )

            module = {
                "id": module_id,
                "name": metadata.get(
                    "name",
                    blueprint_name
                ),
                "status": "running",
                "blueprint": blueprint_name,
                "runtime": runtime.get(
                    "using",
                    "unknown"
                ),
                "started": int(
                    time.time()
                )
            }

            os.makedirs(
                "data/modules",
                exist_ok=True
            )

            module_file = os.path.join(
                "data",
                "modules",
                f"{module_id}.json"
            )

            with open(
                module_file,
                "w",
                encoding="utf-8"
            ) as file:

                json.dump(
                    module,
                    file,
                    indent=4
                )

            print()

            print(
                f"Igniting blueprint: {blueprint_name}"
            )

            print(
                f"Module ID: {module_id}"
            )

            print(
                f"Runtime: {module['runtime']}"
            )

            print(
                f"Status: {module['status']}"
            )

            print()

        except Exception as error:

            print(
                f"Ignite failed: {error}"
            )