class DevLexer:

    KEYWORDS = {
        "development",
        "using",
        "workspace",
        "service",
        "port",
        "watch",
        "env",
        "startup",
        "shell"
    }

    def lex(self, tokens):

        result = []

        for token in tokens:

            if token in self.KEYWORDS:

                result.append(
                    ("KEYWORD", token)
                )

            elif token == ":|":

                result.append(
                    ("BLOCK_START", token)
                )

            elif token == "|:":

                result.append(
                    ("BLOCK_END", token)
                )

            else:

                result.append(
                    ("VALUE", token)
                )

        return result