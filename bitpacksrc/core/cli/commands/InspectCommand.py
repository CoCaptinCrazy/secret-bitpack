import json
import os


class InspectCommand:

    def execute(self, args):

        if not args:
            print("Usage: bitpack inspect <blueprint>")
            return

        blueprint = args[0]

        path = os.path.join(
            "data",
            "blueprints",
            f"{blueprint}.bp"
        )

        if not os.path.exists(path):
            print("Blueprint not found.")
            return

        with open(
            path,
            "r",
            encoding="utf-8"
        ) as file:

            data = json.load(file)

        print(
            json.dumps(
                data,
                indent=4
            )
        )