class DashboardManager:

    def __init__(self):
        self.dashboards = {}

    def register(
        self,
        dashboard
    ):
        self.dashboards[
            dashboard.name
        ] = dashboard