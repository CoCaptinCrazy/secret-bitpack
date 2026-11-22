from bitpack.networking.mesh.MeshManager import MeshManager


class MeshCommand:

    def execute(self, args):

        manager = MeshManager()

        if not args:
            print(
                manager.list()
            )
            return

        name = args[0]

        mesh = manager.create(name)

        print(
            f"Mesh created: {mesh.name}"
        )