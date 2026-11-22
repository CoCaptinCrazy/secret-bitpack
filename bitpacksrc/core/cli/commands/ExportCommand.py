import shutil
import os


class ExportCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bitpack export "
                "<blueprint>"
            )

            return

        source = os.path.join(
            "data",
            "blueprints",
            f"{args[0]}.bp"
        )

        target = (
            f"{args[0]}.export.bp"
        )

        shutil.copy(
            source,
            target
        )

        print(
            f"Exported to "
            f"{target}"
        )