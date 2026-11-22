from bitpack.core.cli.HelpGenerator import HelpGenerator


class CommandRouter:

    def __init__(self):
        from bitpack.core.cli.CommandRegistry import CommandRegistry
        self.registry = CommandRegistry()

    def route(self, command, args):

        handler = self.registry.get(command)

        if handler:
            handler.execute(args)
            return

        HelpGenerator().unknown(command)