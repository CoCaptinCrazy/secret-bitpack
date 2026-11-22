class BlueprintGC:

    def collect(
        self,
        blueprints,
        referenced
    ):

        garbage = []

        for name in blueprints:

            if name not in referenced:
                garbage.append(
                    name
                )

        return garbage