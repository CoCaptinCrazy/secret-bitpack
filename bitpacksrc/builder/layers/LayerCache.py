class LayerCache:

    def __init__(self, store):
        self.store = store

    def hit(self, layer):

        return self.store.exists(
            layer.id
        )

    def cache(self, layer):

        self.store.save(layer)