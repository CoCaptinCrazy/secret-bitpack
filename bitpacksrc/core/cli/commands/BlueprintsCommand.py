import json
import os
import datetime


class BlueprintsCommand:

    def execute(self, args):

        blueprint_dir = "data/blueprints"

        if not os.path.exists(blueprint_dir):
            print("No blueprints found.")
            return

        files = [
            f
            for f in os.listdir(blueprint_dir)
            if f.endswith(".bp")
        ]

        if not files:
            print("No blueprints found.")
            return

        print()

        print(
            f"{'NAME':20}"
            f"{'VERSION':10}"
            f"{'RUNTIME':12}"
            f"{'LAYERS':8}"
            f"{'VOLUMES':10}"
            f"{'SECRETS':10}"
            f"{'MESHES':8}"
            f"{'CREATED':20}"
            f"{'ID':12}"
        )

        print("-" * 90)

        for file in files:

            path = os.path.join(
                blueprint_dir,
                file
            )

            try:

                with open(
                    path,
                    "r",
                    encoding="utf-8"
                ) as f:

                    bp = json.load(f)

                metadata = bp.get(
                    "metadata",
                    {}
                )

                runtime = bp.get(
                    "runtime",
                    {}
                )

                layers = bp.get(
                    "layers",
                    []
                )

                name = metadata.get(
                    "name",
                    "unknown"
                )

                version = metadata.get(
                    "version",
                    "unknown"
                )

                created = metadata.get(
                    "created",
                    0
                )

                blueprint_id = metadata.get(
                    "id",
                    "unknown"
                )[:12]

                runtime_name = runtime.get(
                    "using",
                    "unknown"
                )

                layer_count = len(
                    layers
                )

                storage = bp.get(
                    "storage",
                    {}
                )

                security = bp.get(
                    "security",
                    {}
                )

                networking = bp.get(
                    "networking",
                    {}
                )

                volumes = len(
                    storage.get(
                        "volumes",
                        []
                    )
                )

                secrets = len(
                    security.get(
                        "secrets",
                        []
                    )
                )

                meshes = len(
                    networking.get(
                        "meshes",
                        []
                    )
                )

                if created:
                    created_text = (
                        datetime.datetime
                        .fromtimestamp(created)
                        .strftime(
                            "%Y-%m-%d %H:%M"
                        )
                    )
                else:
                    created_text = "-"

                print(
                    f"{name:20}"
                    f"{version:10}"
                    f"{runtime_name:12}"
                    f"{layer_count:<8}"
                    f"{volumes:<10}"
                    f"{secrets:<10}"
                    f"{meshes:<8}"
                    f"{created_text:20}"
                    f"{blueprint_id:12}"
                )

            except Exception as error:

                print(
                    f"{file:20}"
                    f"{'ERROR':10}"
                    f"{'-':12}"
                    f"{'-':8}"
                    f"{'-':20}"
                    f"{'-':12}"
                )

        print()