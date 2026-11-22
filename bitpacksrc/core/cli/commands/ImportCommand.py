import shutil
import os


class ImportCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bitpack import "
                "<file.bp>"
            )

            return

        source = args[0]

        destination = os.path.join(
            "data",
            "blueprints",
            os.path.basename(
                source
            )
        )

        shutil.copy(
            source,
            destination
        )

        print(
            f"Imported: "
            f"{destination}"
        )