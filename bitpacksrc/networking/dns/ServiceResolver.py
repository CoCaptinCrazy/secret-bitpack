class ServiceResolver:

    def resolve(
        self,
        mesh,
        service_name
    ):

        if service_name not in mesh.modules:
            return None

        return mesh.modules[
            service_name
        ]