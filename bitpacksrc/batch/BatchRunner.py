from bitpack.batch.runtime.BatchEngine import BatchEngine


class BatchRunner:

    def __init__(self):

        self.engine = BatchEngine()

    def run(
        self,
        plan
    ):

        print(
            f"[PIPELINE] {plan.name}"
        )

        return self.engine.execute(
            plan
        )