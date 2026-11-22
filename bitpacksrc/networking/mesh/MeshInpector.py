class MeshInspector:

    def inspect(self, mesh):

        return {
            "id": mesh.id,
            "name": mesh.name,
            "modules": list(
                mesh.modules.keys()
            )
        }