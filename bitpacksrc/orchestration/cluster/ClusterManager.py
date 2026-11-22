import os
import json
import time

from bitpack.orchestration.cluster.Cluster import Cluster


class ClusterManager:

    def __init__(self):

        self.path = "data/clusters"

        os.makedirs(
            self.path,
            exist_ok=True
        )

    def create(
        self,
        name
    ):

        cluster = Cluster(name)

        self._save(
            cluster,
            {
                "kind": "Cluster",
                "apiVersion": "v1",

                "metadata": {
                    "name": name,
                    "created": int(
                        time.time()
                    ),
                    "bitpackVersion": "0.1.17"
                },

                "status": {
                    "phase": "Ready",
                    "healthy": True
                },

                "nodes": [],

                "assemblies": [],

                "modules": [],

                "meshes": [],

                "snapshots": [],

                "events": [
                    {
                        "type": "ClusterCreated",
                        "message":
                            f"Cluster {name} created",
                        "timestamp":
                            int(time.time())
                    }
                ],

                "metrics": {
                    "nodes": 0,
                    "assemblies": 0,
                    "modules": 0,
                    "meshes": 0,
                    "snapshots": 0
                }
            }
        )

        return cluster

    def get(
        self,
        name
    ):

        file = os.path.join(
            self.path,
            f"{name}.json"
        )

        if not os.path.exists(file):
            return None

        with open(
            file,
            "r",
            encoding="utf-8"
        ) as f:

            return json.load(f)

    def delete(
        self,
        name
    ):

        file = os.path.join(
            self.path,
            f"{name}.json"
        )

        if os.path.exists(file):

            os.remove(file)

            return True

        return False

    def list(self):

        result = []

        for file in os.listdir(self.path):

            if file.endswith(".json"):

                result.append(
                    Cluster(
                        file[:-5]
                    )
                )

        return sorted(
            result,
            key=lambda c: c.name
        )


    def add_node(
        self,
        cluster_name,
        node
    ):

        data = self.get(
            cluster_name
        )

        if not data:
            return False

        data["nodes"].append(node)

        data["metrics"]["nodes"] = len(
            data["nodes"]
        )

        self._event(
            data,
            "NodeAdded",
            f"Added node {node}"
        )

        self._save_raw(data)

        return True

    def add_module(
        self,
        cluster_name,
        module
    ):

        data = self.get(
            cluster_name
        )

        if not data:
            return False

        if module not in data["modules"]:

            data["modules"].append(
                module
            )

        data["metrics"]["modules"] = len(
            data["modules"]
        )

        self._event(
            data,
            "ModuleStarted",
            f"Started module {module}"
        )

        self._save_raw(data)

        return True

    def add_assembly(
        self,
        cluster_name,
        assembly
    ):

        data = self.get(
            cluster_name
        )

        if not data:
            return False

        if assembly not in data["assemblies"]:

            data["assemblies"].append(
                assembly
            )

        data["metrics"][
            "assemblies"
        ] = len(
            data["assemblies"]
        )

        self._event(
            data,
            "AssemblyDeployed",
            f"Deployed assembly {assembly}"
        )

        self._save_raw(data)

        return True

    def add_mesh(
        self,
        cluster_name,
        mesh
    ):

        data = self.get(
            cluster_name
        )

        if not data:
            return False

        if mesh not in data["meshes"]:

            data["meshes"].append(
                mesh
            )

        data["metrics"]["meshes"] = len(
            data["meshes"]
        )

        self._event(
            data,
            "MeshCreated",
            f"Created mesh {mesh}"
        )

        self._save_raw(data)

        return True

    def add_snapshot(
        self,
        cluster_name,
        snapshot
    ):

        data = self.get(
            cluster_name
        )

        if not data:
            return False

        data["snapshots"].append(
            snapshot
        )

        data["metrics"][
            "snapshots"
        ] = len(
            data["snapshots"]
        )

        self._event(
            data,
            "SnapshotCreated",
            f"Created snapshot {snapshot}"
        )

        self._save_raw(data)

        return True

    def _event(
        self,
        data,
        event_type,
        message
    ):

        data["events"].append(
            {
                "type": event_type,
                "message": message,
                "timestamp": int(
                    time.time()
                )
            }
        )

    def _save(
        self,
        cluster,
        data
    ):

        file = os.path.join(
            self.path,
            f"{cluster.name}.json"
        )

        with open(
            file,
            "w",
            encoding="utf-8"
        ) as f:

            json.dump(
                data,
                f,
                indent=4
            )

    def _save_raw(
        self,
        data
    ):

        file = os.path.join(
            self.path,
            f"{data['metadata']['name']}.json"
        )

        with open(
            file,
            "w",
            encoding="utf-8"
        ) as f:

            json.dump(
                data,
                f,
                indent=4
            )