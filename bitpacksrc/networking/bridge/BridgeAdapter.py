class BridgeAdapter:

    def map(
        self,
        bridge
    ):

        return {
            "host":
                bridge.host_port,

            "target":
                bridge.target_port
        }