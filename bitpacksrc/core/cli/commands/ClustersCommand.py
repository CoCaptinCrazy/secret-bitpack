from bitpack.orchestration.cluster.ClusterManager import ClusterManager


class ClustersCommand:

    def execute(self, args):

        manager = ClusterManager()

        clusters = manager.list()

        if not clusters:

            print(
                "No clusters."
            )

            return

        for cluster in clusters:

            print(
                "- " + cluster.name
            )