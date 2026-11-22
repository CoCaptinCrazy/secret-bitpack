class ASTNode:

    def __init__(
        self,
        node_type,
        value=None
    ):
        self.type = node_type
        self.value = value
        self.children = []

    def add(self, node):
        self.children.append(node)

    def __repr__(self):
        return (
            f"ASTNode("
            f"type={self.type}, "
            f"value={self.value})"
        )