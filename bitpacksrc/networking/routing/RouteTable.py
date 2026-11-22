class RouteTable:

    def __init__(self):
        self.routes = {}

    def add(
        self,
        source,
        destination
    ):
        self.routes[source] = destination

    def lookup(
        self,
        source
    ):
        return self.routes.get(
            source
        )