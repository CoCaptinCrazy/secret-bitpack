from bitpack.core.Bitpack import Bitpack


class StopCommand:

    def execute(self, args):

        if not args:
            print(
                "Usage: bitpack stop <module>"
            )
            return

        module_id = args[0]

        bp = Bitpack()

        try:

            bp.runtime.stop(
                module_id
            )

            print(
                f"Stopped: {module_id}"
            )

        except Exception as e:

            print(
                f"Stop failed: {e}"
            )