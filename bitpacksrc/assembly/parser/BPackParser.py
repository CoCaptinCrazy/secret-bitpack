from bitpack.assembly.parser.AssemblyGraph import AssemblyNode
from bitpack.assembly.parser.BPackTokenizer import BPackTokenizer


class BPackParser:

    def parse(self, source):

        tokenizer = BPackTokenizer()

        tokens = tokenizer.tokenize(source)

        root = AssemblyNode(
            "root",
            "root"
        )

        stack = [root]

        i = 0

        while i < len(tokens):

            token = tokens[i]

            if (
                token == "assembly"
                and i + 1 < len(tokens)
            ):

                node = AssemblyNode(
                    "assembly",
                    tokens[i + 1]
                )

                stack[-1].add(node)
                stack.append(node)

                i += 2
                continue

            if (
                token == "mesh"
                and i + 1 < len(tokens)
            ):

                node = AssemblyNode(
                    "mesh",
                    tokens[i + 1]
                )

                stack[-1].add(node)
                stack.append(node)

                i += 2
                continue

            if (
                token == "module"
                and i + 1 < len(tokens)
            ):

                node = AssemblyNode(
                    "module",
                    tokens[i + 1]
                )

                stack[-1].add(node)
                stack.append(node)

                i += 2
                continue

            if token == "}":

                if len(stack) > 1:
                    stack.pop()

            i += 1

        return root