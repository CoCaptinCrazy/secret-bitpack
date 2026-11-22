from bitpack.builder.parser.Token import Token


class Lexer:

    KEYWORDS = {
        "blueprint",
        "using",
        "bring",
        "env",
        "assemble",
        "verify",
        "ignite",
        "health",
        "secret"
    }

    def lex(self, raw_tokens):

        result = []

        for token in raw_tokens:

            if token in self.KEYWORDS:

                result.append(
                    Token(
                        "KEYWORD",
                        token
                    )
                )

            elif token == "<|":

                result.append(
                    Token(
                        "BLOCK_START",
                        token
                    )
                )

            elif token == "|>":

                result.append(
                    Token(
                        "BLOCK_END",
                        token
                    )
                )

            else:

                result.append(
                    Token(
                        "VALUE",
                        token
                    )
                )

        return result