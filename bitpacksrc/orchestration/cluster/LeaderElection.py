class LeaderElection:

    def elect(
        self,
        nodes
    ):

        if not nodes:
            return None

        return min(
            nodes,
            key=lambda node: node.id
        )