from bitpack.builder.parser.AST import ASTNode


class Parser:

    def parse(self, tokens):

        root = ASTNode("ROOT")

        i = 0

        while i < len(tokens):

            token = tokens[i]

            if token.type == "KEYWORD":

                value = None

                # grab next token as value if present
                if (
                    i + 1 < len(tokens)
                    and tokens[i + 1].type == "VALUE"
                ):
                    value = tokens[i + 1].value

                node = ASTNode(
                    token.value.upper(),
                    value
                )

                root.add(node)

            i += 1

        return root