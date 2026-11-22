from bitpack.monitoring.metrics.MetricsRegistry import MetricsRegistry


class MetricsCommand:

    def execute(self, args):

        registry = (
            MetricsRegistry()
        )

        metrics = (
            registry.all()
        )

        if not metrics:

            print(
                "No metrics."
            )

            return

        for name in metrics:

            metric = metrics[name]

            print(
                f"{metric.name}: "
                f"{metric.value}"
            )