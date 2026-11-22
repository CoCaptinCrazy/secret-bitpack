class ScaleCommand:

    def execute(self, args):

        if len(args) < 2:

            print(
                "Usage: bitpack scale "
                "<service> <replicas>"
            )

            return

        service = args[0]

        replicas = int(
            args[1]
        )

        print(
            f"Scaling "
            f"{service} "
            f"to "
            f"{replicas}"
        )