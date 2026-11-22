from bitpack.dev.parser.DevParser import DevParser
from bitpack.dev.runtime.DevRuntime import DevRuntime


class DevCommand:

    def execute(self, args):

        file = (
            args[0]
            if args
            else "project.bpdev"
        )

        with open(
            file,
            "r",
            encoding="utf-8"
        ) as f:

            source = f.read()

        parser = DevParser()

        environment = parser.parse(
            source
        )

        runtime = DevRuntime()

        runtime.launch(
            environment
        )