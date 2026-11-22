class AssemblyNode:

    def __init__(
        self,
        node_type,
        name
    ):
        self.type = node_type
        self.name = name

        self.children = []

    def add(self, node):
        self.children.append(node)