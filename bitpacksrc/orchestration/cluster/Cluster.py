class Cluster:

    def __init__(
        self,
        name
    ):
        self.name = name

        self.nodes = []

        self.assemblies = []

        self.modules = []

        self.status = "ready"

    def add_node(
        self,
        node
    ):
        self.nodes.append(node)

    def remove_node(
        self,
        node_id
    ):
        self.nodes = [
            node
            for node in self.nodes
            if node.id != node_id
        ]

    def add_assembly(
        self,
        assembly
    ):

        if assembly not in self.assemblies:

            self.assemblies.append(
                assembly
            )

    def remove_assembly(
        self,
        assembly
    ):

        self.assemblies = [
            item
            for item in self.assemblies
            if item != assembly
        ]

    def add_module(
        self,
        module
    ):

        if module not in self.modules:

            self.modules.append(
                module
            )

    def remove_module(
        self,
        module
    ):

        self.modules = [
            item
            for item in self.modules
            if item != module
        ]

    def node_count(
        self
    ):

        return len(
            self.nodes
        )

    def module_count(
        self
    ):

        return len(
            self.modules
        )

    def assembly_count(
        self
    ):

        return len(
            self.assemblies
        )