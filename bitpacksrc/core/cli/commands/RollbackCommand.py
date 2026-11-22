class RollbackCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack rollback "
                "<assembly>"
            )

            return

        print(
            f"Rolling back "
            f"{args[0]}"
        )