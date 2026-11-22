from bitpack.monitoring.alerts.AlertManager import AlertManager
from bitpack.monitoring.alerts.AlertPolicy import AlertPolicy


class AlertsCommand:

    def execute(self, args):

        manager = AlertManager()

        policy = AlertPolicy()

        alerts = manager.evaluate(
            {
                "cpu": 100,
                "memory": 20
            },
            policy
        )

        if not alerts:

            print(
                "No alerts."
            )

            return

        for alert in alerts:

            print(
                f"[{alert.severity}] "
                f"{alert.message}"
            )