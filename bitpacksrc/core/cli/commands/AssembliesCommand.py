import os


class AssembliesCommand:

    def execute(self, args):

        path = os.path.join(
            "data",
            "assemblies"
        )

        if not os.path.exists(path):

            print(
                "No assemblies found."
            )

            return

        assemblies = [

            item

            for item in os.listdir(path)

            if os.path.isdir(
                os.path.join(
                    path,
                    item
                )
            )

        ]

        if not assemblies:

            print(
                "No assemblies found."
            )

            return

        print()
        print("ASSEMBLIES")
        print(
            "-" * 30
        )

        for assembly in sorted(
            assemblies
        ):

            print(assembly)

        print()