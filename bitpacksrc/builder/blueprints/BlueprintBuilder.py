from bitpack.builder.blueprints.Blueprint import Blueprint
from bitpack.builder.blueprints.BlueprintVerifier import BlueprintVerifier
from bitpack.builder.blueprints.BlueprintManifest import BlueprintManifest


class BlueprintBuilder:

    def __init__(self):
        self.verifier = BlueprintVerifier()
        self.manifest = BlueprintManifest()

    def build(
        self,
        name,
        version,
        layers,
        metadata
    ):

        blueprint = Blueprint(
            name,
            version
        )

        for layer in layers:
            blueprint.add_layer(layer)

        blueprint.set_metadata(
            metadata
        )

        if not self.verifier.verify(
            blueprint
        ):
            raise ValueError(
                "Invalid blueprint"
            )

        return blueprint