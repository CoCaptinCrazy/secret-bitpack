class BPackLexer:

    KEYWORDS = {
        "assembly",
        "mesh",
        "module",
        "blueprint",
        "depends",
        "volume",
        "firewall",
        "autoscale",
        "policy"
    }

    def lex(self, tokens):

        result = []

        for token in tokens:

            if token in self.KEYWORDS:
                result.append(("KEYWORD", token))

            elif token in ["{", "}"]:
                result.append(("BRACE", token))

            else:
                result.append(("VALUE", token))

        return result