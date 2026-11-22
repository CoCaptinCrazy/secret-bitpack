from bitpack.monitoring.profiling.Profiler import Profiler


class ProfileCommand:

    def execute(self, args):

        profiler = Profiler()

        trace = profiler.begin()

        print(
            "Profile started."
        )

        elapsed = trace.stop()

        print(
            f"Elapsed: "
            f"{elapsed:.6f}s"
        )