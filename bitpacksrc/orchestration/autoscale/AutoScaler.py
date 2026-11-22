class AutoScaler:

    def evaluate(
        self,
        replicas,
        metrics,
        policy
    ):

        if (
            metrics.cpu >
            policy.cpu_threshold
        ):
            return min(
                replicas + 1,
                policy.max_replicas
            )

        return replicas