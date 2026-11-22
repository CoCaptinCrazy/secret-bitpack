from bitpack.monitoring.alerts.AlertManager import AlertManager
from bitpack.monitoring.alerts.AlertPolicy import AlertPolicy


class AlertCommand:

    def execute(self, args):

        manager = AlertManager()

        policy = AlertPolicy()

        alerts = manager.evaluate(
            {
                "cpu": 95,
                "memory": 40
            },
            policy
        )

        for alert in alerts:

            print(
                f"[{alert.severity}] "
                f"{alert.message}"
            )