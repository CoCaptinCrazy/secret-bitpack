class RuntimeConfig:
    def __init__(self):
        self.debug = False
        self.verbose = False
        self.offline = False
        self.dry_run = False

    def load(self):
        pass