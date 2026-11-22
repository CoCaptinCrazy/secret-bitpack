from bitpack.runtime.health.Probe import Probe


class HealthCheck:

    def __init__(self):
        self.probe = Probe()

    def run(
        self,
        module
    ):

        return self.probe.check(
            module
        )