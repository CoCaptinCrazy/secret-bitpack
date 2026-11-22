from bitpack.monitoring.alerts.Alert import Alert


class AlertManager:

    def evaluate(
        self,
        metrics,
        policy
    ):

        alerts = []

        if metrics["cpu"] > policy.cpu_limit:

            alerts.append(
                Alert(
                    "warning",
                    "CPU limit exceeded"
                )
            )

        return alerts