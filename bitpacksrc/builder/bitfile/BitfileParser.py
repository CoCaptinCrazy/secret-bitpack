from bitpack.builder.parser.Tokenizer import Tokenizer
from bitpack.builder.parser.Lexer import Lexer
from bitpack.builder.parser.Parser import Parser
from bitpack.builder.parser.Validator import Validator


class BitfileParser:

    def __init__(self):
        self.tokenizer = Tokenizer()
        self.lexer = Lexer()
        self.parser = Parser()
        self.validator = Validator()

    def parse_file(self, path):

        with open(
            path,
            "r",
            encoding="utf-8"
        ) as file:

            source = file.read()

        return self.parse(source)

    def parse(self, source):

        raw_tokens = (
            self.tokenizer.tokenize(source)
        )

        tokens = (
            self.lexer.lex(raw_tokens)
        )

        ast = (
            self.parser.parse(tokens)
        )

        self.validator.validate(ast)

        return ast