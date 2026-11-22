import json
import os


class VerifyCommand:

    def execute(self, args):

        if not args:
            print("Usage: bitpack verify <blueprint>")
            return

        path = os.path.join(
            "data",
            "blueprints",
            f"{args[0]}.bp"
        )

        if not os.path.exists(path):
            print("Blueprint not found.")
            return

        with open(path) as file:
            blueprint = json.load(file)

        required = [
            "metadata",
            "runtime",
            "execution"
        ]

        missing = []

        for item in required:

            if item not in blueprint:
                missing.append(item)

        if missing:

            print(
                "INVALID"
            )

            print(
                "Missing:"
            )

            for item in missing:
                print(f" - {item}")

            return

        print("VALID")