from bitpack.orchestration.cluster.ClusterManager import ClusterManager


class ClusterCommand:

    def execute(self, args):

        if not args:

            print(
                "Usage: bitpack cluster <action> [options]"
            )

            return

        action = args[0]

        manager = ClusterManager()

        if action == "create":

            if len(args) < 2:

                print(
                    "Usage: bitpack cluster create <name>"
                )

                return

            cluster = manager.create(
                args[1]
            )

            print(
                f"Cluster created: {cluster.name}"
            )

        elif action == "delete" or "rm":

            if len(args) < 2:

                print(
                    "Usage: bitpack cluster delete <name>"
                )

                return

            manager.delete(
                args[1]
            )

            print(
                f"Cluster deleted: {args[1]}"
            )

        elif action == "inspect":

            if len(args) < 2:

                print(
                    "Usage: bitpack cluster inspect <name>"
                )

                return

            cluster = manager.load(
                args[1]
            )

            if not cluster:

                print(
                    "Cluster not found."
                )

                return

            print()
            print(
                f"Name: {cluster['name']}"
            )
            print(
                f"Nodes: {cluster['nodes']}"
            )
            print(
                f"Status: {cluster['status']}"
            )
            print()

        else:

            print(
                f"Unknown cluster action: {action}"
            )