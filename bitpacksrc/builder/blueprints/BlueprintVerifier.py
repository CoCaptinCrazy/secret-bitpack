class BlueprintVerifier:

    def verify(
        self,
        blueprint
    ):

        if not blueprint.name:
            return False

        if not blueprint.layers:
            return False

        return True