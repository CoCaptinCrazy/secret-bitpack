from bitpack.orchestration.cluster.ClusterManager import ClusterManager


class NodesCommand:

    def execute(self, args):

        manager = ClusterManager()

        clusters = manager.list()

        if not clusters:

            print(
                "No nodes."
            )

            return

        for cluster in clusters:

            print()

            print(
                f"Cluster: "
                f"{cluster.name}"
            )

            for node in cluster.nodes:

                print(
                    f" - {node.name}"
                )