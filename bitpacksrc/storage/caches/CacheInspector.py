class CacheInspector:

    ROOT = "./data/caches"

    def list(self):

        return [
            file.stem
            for file
            in self.ROOT.glob("*.json")
        ]