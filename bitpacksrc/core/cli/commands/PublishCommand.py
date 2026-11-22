from bitpack.registry.push.PushManager import PushManager


class PublishCommand:

    def execute(self, args):

        if not args:
            print("Usage: bitpack publish <blueprint>")
            return

        blueprint = args[0]

        try:

            push = PushManager()

            result = push.publish(
                blueprint
            )

            print(
                f"Published: {blueprint}"
            )

            print(result)

        except Exception as e:

            print(
                f"Publish failed: {e}"
            )