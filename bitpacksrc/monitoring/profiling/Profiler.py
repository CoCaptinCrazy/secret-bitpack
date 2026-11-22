from bitpack.monitoring.profiling.Trace import Trace


class Profiler:

    def begin(self):
        return Trace()