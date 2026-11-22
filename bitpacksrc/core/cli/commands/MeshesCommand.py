import os


class MeshesCommand:

    def execute(self, args):

        root = (
            "data/meshes"
        )

        if not os.path.exists(root):

            print(
                "No meshes found."
            )

            return

        meshes = os.listdir(root)

        if not meshes:

            print(
                "No meshes found."
            )

            return

        for mesh in meshes:

            print(mesh)