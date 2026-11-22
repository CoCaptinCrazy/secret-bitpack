import os
import json
import time
import hashlib

from bitpack.ignore.IgnoreEngine import IgnoreEngine


class ForgeCommand:

    def execute(self, args):

        bitfile = (
            args[0]
            if args
            else "Bitfile"
        )

        if not os.path.exists(bitfile):

            print(
                f"Bitfile not found: {bitfile}"
            )

            return

        try:

            with open(
                bitfile,
                "r",
                encoding="utf-8"
            ) as file:

                source = file.read()

            ignore = IgnoreEngine()

            forge_context = ignore.collect(".")

            blueprint_name = self._extract_name(
                source
            )

            if not blueprint_name:

                blueprint_name = "unnamed"

            blueprint_id = hashlib.sha256(
                source.encode()
            ).hexdigest()

            blueprint = {

                "kind": "Blueprint",

                "apiVersion": "v1",

                "metadata": {

                    "id": blueprint_id,

                    "name": blueprint_name,

                    "version": "1.0",

                    "created": int(
                        time.time()
                    ),

                    "author": "unknown",

                    "description": "",

                    "bitpackVersion": "0.1.0"
                },

                "runtime": {

                    "using": "unknown",

                    "environment": {},

                    "health": []
                },

                "execution": {

                    "assemble": [],

                    "verify": [],

                    "ignite": []
                },

                "storage": {

                    "volumes": [],

                    "caches": [],

                    "snapshots": True
                },

                "networking": {

                    "meshes": [],

                    "bridges": [],

                    "publishedPorts": []
                },

                "security": {

                    "secrets": [],

                    "signatures": [],

                    "policies": []
                },

                "monitoring": {

                    "metrics": True,

                    "alerts": True,

                    "profiling": True
                },

                "assembly": {

                    "compatible": True,

                    "clusterReady": True,

                    "autoscaleReady": True
                },

                "layers": [],

                "history": [

                    {
                        "action": "forge",

                        "timestamp": int(
                            time.time()
                        )
                    }

                ],

                "source": {

                    "bitfile": source,

                    "context": forge_context
                }
            }

            os.makedirs(
                "data/blueprints",
                exist_ok=True
            )

            output = os.path.join(

                "data",

                "blueprints",

                f"{blueprint_name}.bp"
            )

            with open(

                output,

                "w",

                encoding="utf-8"

            ) as file:

                json.dump(

                    blueprint,

                    file,

                    indent=4

                )

            print()
            print(
                f"Blueprint created: {output}"
            )
            print(
                f"Included files: {len(forge_context)}"
            )

        except Exception as error:

            print(
                f"Forge failed: {error}"
            )

    def _extract_name(
        self,
        source
    ):

        for line in source.splitlines():

            line = line.strip()

            if line.startswith(
                "blueprint"
            ):

                parts = line.split(
                    '"'
                )

                if len(parts) >= 2:

                    return parts[1]

                pieces = line.split()

                if len(pieces) >= 2:

                    return pieces[1]

        return None