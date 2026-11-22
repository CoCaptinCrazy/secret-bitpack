class BlueprintManifest:

    def generate(
        self,
        blueprint
    ):
        return {
            "name":
                blueprint.name,

            "version":
                blueprint.version,

            "layers":
                [
                    layer.id
                    for layer
                    in blueprint.layers
                ],

            "metadata":
                blueprint.metadata
        }