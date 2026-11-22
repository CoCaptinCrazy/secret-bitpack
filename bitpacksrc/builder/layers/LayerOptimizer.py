class LayerOptimizer:

    def optimize(
        self,
        layers
    ):

        optimized = []

        seen = set()

        for layer in layers:

            if layer.id in seen:
                continue

            seen.add(
                layer.id
            )

            optimized.append(
                layer
            )

        return optimized