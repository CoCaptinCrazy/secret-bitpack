class ServiceDiscovery:

    def resolve(
        self,
        services,
        name
    ):

        return services.get(
            name
        )