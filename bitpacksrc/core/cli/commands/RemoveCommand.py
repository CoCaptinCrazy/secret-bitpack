from bitpack.core.Bitpack import Bitpack


class RemoveCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack rm <module>"
            )
            return

        module_id = args[0]

        bp = Bitpack()

        try:

            bp.state.remove(
                module_id
            )

            print(
                f"Removed: {module_id}"
            )

        except Exception as e:

            print(
                f"Remove failed: {e}"
            )