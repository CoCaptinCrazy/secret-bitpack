from bitpack.networking.bridge.Bridge import Bridge


class BridgeManager:

    def __init__(self):
        self.bridges = {}

    def create(
        self,
        name,
        host_port,
        target_port
    ):

        bridge = Bridge(
            name,
            host_port,
            target_port
        )

        self.bridges[
            bridge.id
        ] = bridge

        return bridge

    def remove(
        self,
        bridge_id
    ):

        self.bridges.pop(
            bridge_id,
            None
        )

    def list(self):

        return list(
            self.bridges.values()
        )