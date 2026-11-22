from bitpack.runtime.health.HealthCheck import HealthCheck


class Watchdog:

    def __init__(self):
        self.health = HealthCheck()

    def watch(
        self,
        module
    ):

        healthy = self.health.run(
            module
        )

        if not healthy:
            self.restart(module)

    def restart(
        self,
        module
    ):

        print(
            f"Restarting {module.name}"
        )