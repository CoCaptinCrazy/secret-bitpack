class MetricsRegistry:

    def __init__(self):
        self.metrics = {}

    def register(
        self,
        metric
    ):
        self.metrics[
            metric.name
        ] = metric

    def get(
        self,
        name
    ):
        return self.metrics.get(name)

    def all(self):
        return self.metrics