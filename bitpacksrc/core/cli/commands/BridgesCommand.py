from bitpack.networking.bridge.BridgeManager import BridgeManager


class BridgesCommand:

    def execute(self, args):

        manager = BridgeManager()

        bridges = manager.list()

        if not bridges:

            print(
                "No bridges."
            )

            return

        for bridge in bridges:

            print(
                f"{bridge.name} "
                f"{bridge.host_port}"
                f"->{bridge.target_port}"
            )