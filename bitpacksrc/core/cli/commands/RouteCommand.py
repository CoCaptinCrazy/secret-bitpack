from bitpack.networking.routing.RouteTable import RouteTable


class RouteCommand:

    def execute(self, args):

        table = RouteTable()

        print(
            table.routes
        )