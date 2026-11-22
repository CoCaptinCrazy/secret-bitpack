class RolloutCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack rollout "
                "<assembly>"
            )

            return

        print(
            f"Rolling out "
            f"{args[0]}"
        )