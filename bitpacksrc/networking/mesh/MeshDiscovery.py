class MeshDiscovery:

    def resolve(
        self,
        mesh,
        module_name
    ):

        service = mesh.modules.get(
            module_name
        )

        if not service:
            return None

        return service