class Router:

    def __init__(
        self,
        table
    ):
        self.table = table

    def route(
        self,
        destination
    ):
        return self.table.lookup(
            destination
        )