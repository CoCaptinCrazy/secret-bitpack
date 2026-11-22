class TrendAnalyzer:

    def trend(
        self,
        values
    ):

        if len(values) < 2:
            return "stable"

        if values[-1] > values[0]:
            return "up"

        if values[-1] < values[0]:
            return "down"

        return "stable"