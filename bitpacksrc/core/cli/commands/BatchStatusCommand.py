import os


class BatchStatusCommand:

    def execute(self, args):

        root = (
            "data/batch"
        )

        if not os.path.exists(root):

            print(
                "No batch jobs."
            )

            return

        jobs = os.listdir(root)

        if not jobs:

            print(
                "No batch jobs."
            )

            return

        for job in jobs:

            print(job)