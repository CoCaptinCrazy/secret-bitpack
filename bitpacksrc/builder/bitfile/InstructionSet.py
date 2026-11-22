class InstructionSet:

    def __init__(self):
        self.blueprint = None
        self.version = None

        self.using = None

        self.bring = []

        self.env = {}

        self.assemble = []
        self.verify = []
        self.ignite = []

        self.health = []
        self.secrets = []