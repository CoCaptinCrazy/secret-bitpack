import os
import time


class WatchManager:

    def watch(self, paths):

        print(
            "[WATCH] File monitoring enabled"
        )

        for path in paths:

            print(
                f"[WATCH] {path}"
            )

    def snapshot(self, path):

        state = {}

        for root, dirs, files in os.walk(path):

            for file in files:

                full_path = os.path.join(
                    root,
                    file
                )

                try:

                    state[
                        full_path
                    ] = os.path.getmtime(
                        full_path
                    )

                except:

                    pass

        return state

    def monitor(
        self,
        path,
        interval=1
    ):

        previous = self.snapshot(path)

        print(
            f"[WATCH] Monitoring {path}"
        )

        while True:

            current = self.snapshot(path)

            for file in current:

                if (
                    file in previous
                    and current[file]
                    != previous[file]
                ):

                    print(
                        f"[CHANGED] {file}"
                    )

            previous = current

            time.sleep(interval)