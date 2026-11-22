import sys

class ArgumentParser:

    def parse(self):

        if len(sys.argv) < 2:
            return "help", []

        command = sys.argv[1]

        args = sys.argv[2:]

        return command, args