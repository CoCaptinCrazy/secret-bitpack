import time
import os


class WatchCommand:

    def execute(self, args):

        target = (
            args[0]
            if args
            else "."
        )

        print(
            f"Watching {target}"
        )

        known = set()

        while True:

            current = set(
                os.listdir(target)
            )

            added = (
                current - known
            )

            for item in added:

                print(
                    f"[CREATE] "
                    f"{item}"
                )

            known = current

            time.sleep(1)