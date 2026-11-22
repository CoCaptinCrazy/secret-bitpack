import json
import os


class ModulesCommand:

    def execute(self, args):

        modules_dir = "data/modules"

        if not os.path.exists(
            modules_dir
        ):

            print(
                "No modules running."
            )

            return

        files = [

            f

            for f in os.listdir(
                modules_dir
            )

            if f.endswith(".json")
        ]

        if not files:

            print(
                "No modules running."
            )

            return

        print()

        print(
            f"{'ID':15}"
            f"{'NAME':20}"
            f"{'STATUS':12}"
            f"{'RUNTIME':10}"
        )

        print(
            "-" * 60
        )

        for file in files:

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

            print(
                f"{module['id'][:14]:15}"
                f"{module['name']:20}"
                f"{module['status']:12}"
                f"{module['runtime']:10}"
            )

        print()