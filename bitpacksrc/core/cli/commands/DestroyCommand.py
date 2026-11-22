import shutil
import os


class DestroyCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack destroy <assembly>"
            )
            return

        target = os.path.join(
            "data",
            "assemblies",
            args[0]
        )

        if not os.path.exists(
            target
        ):

            print(
                "Assembly not found."
            )

            return

        shutil.rmtree(
            target
        )

        print(
            f"Destroyed "
            f"{args[0]}"
        )