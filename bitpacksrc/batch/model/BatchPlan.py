class BatchPlan:

    def __init__(self):

        # Identity
        self.name = ""
        self.version = "1.0"

        self.author = ""
        
        # Actions
        self.forge = []
        self.verify = []
        self.publish = []

        # Configuration
        self.parallel = False
        self.fail_fast = True

        # Reporting
        self.logs = True
        self.metrics = False

        # Dependencies
        self.requires = []

        # Tags
        self.tags = []

        # History
        self.created = None
        self.updated = None

        # Runtime State
        self.status = "pending"

        # Future Support
        self.variables = {}
        self.environment = {}