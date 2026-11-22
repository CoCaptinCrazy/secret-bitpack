class Probe:

    def check(
        self,
        module
    ):

        if module.state.is_running():
            return True

        return False