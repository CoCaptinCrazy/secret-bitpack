from bitpack.networking.mesh.Mesh import Mesh


class MeshManager:

    def __init__(self):
        self.meshes = {}

    def create(self, name):

        mesh = Mesh(name)

        self.meshes[name] = mesh

        return mesh

    def get(self, name):

        return self.meshes.get(name)

    def delete(self, name):

        self.meshes.pop(
            name,
            None
        )

    def list(self):

        return list(
            self.meshes.keys()
        )