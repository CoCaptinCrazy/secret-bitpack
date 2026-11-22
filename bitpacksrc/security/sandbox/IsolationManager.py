class IsolationManager:

    def isolate(
        self,
        module,
        profile
    ):
        module.sandbox = profile