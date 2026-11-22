from bitpack.orchestration.scheduler.Placement import Placement


class Scheduler:

    def __init__(self):
        self.placement = Placement()

    def schedule(
        self,
        module,
        cluster
    ):

        return self.placement.select(
            cluster.nodes
        )