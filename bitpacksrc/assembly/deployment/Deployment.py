import os
import json


class Deployment:

    def deploy(self, root):

        print(f"Deploying assembly: {root.name}")

        self.walk(root)

        return True

    def walk(self, node):

        if node.type != "root":
            print(f"{node.type}: {node.name}")

        if node.type == "assembly":

            os.makedirs(
                f"data/assemblies/{node.name}",
                exist_ok=True
            )

        elif node.type == "mesh":

            os.makedirs(
                f"data/meshes/{node.name}",
                exist_ok=True
            )

        elif node.type == "module":

            module_path = (
                f"data/modules/{node.name}"
            )

            os.makedirs(
                module_path,
                exist_ok=True
            )

            with open(
                f"{module_path}/module.json",
                "w"
            ) as f:

                json.dump(
                    {
                        "name": node.name,
                        "status": "running"
                    },
                    f,
                    indent=4
                )

        for child in node.children:
            self.walk(child)