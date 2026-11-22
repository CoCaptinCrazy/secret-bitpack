from bitpack.runtime.health.HealthCheck import HealthCheck


class HealthCommand:

    def execute(self, args):

        health = HealthCheck()

        if not args:
            print("Usage: bitpack health <module>")
            return

        result = health.run(args[0])

        print(
            f"Health status: {result}"
        )