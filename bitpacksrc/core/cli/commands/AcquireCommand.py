from bitpack.registry.pull.PullManager import PullManager


class AcquireCommand:

    def execute(self, args):

        if not args:
            print("Usage: bitpack acquire <blueprint>")
            return

        blueprint = args[0]

        try:

            pull = PullManager()

            result = pull.acquire(
                blueprint
            )

            print(
                f"Acquired: {blueprint}"
            )

            print(result)

        except Exception as e:

            print(
                f"Acquire failed: {e}"
            )