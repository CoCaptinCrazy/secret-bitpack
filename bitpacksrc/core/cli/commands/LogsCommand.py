import os


class LogsCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bp logs <module>"
            )
            return

        logfile = os.path.join(
            "data",
            "logs",
            "modules",
            f"{args[0]}.log"
        )

        if not os.path.exists(
            logfile
        ):

            print(
                "No logs found."
            )

            return

        with open(logfile) as file:

            print(
                file.read()
            )