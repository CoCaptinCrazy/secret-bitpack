from bitpack.runtime.execution.Executor import Executor
from bitpack.runtime.execution.Environment import Environment


class Launcher:

    def __init__(self):
        self.executor = Executor()

    def launch(
        self,
        module,
        blueprint
    ):

        env = Environment()

        for key, value in blueprint.env.items():
            env.set(
                key,
                value
            )

        command = "\n".join(
            blueprint.ignite
        )

        return self.executor.execute(
            module,
            command,
            env
        )
