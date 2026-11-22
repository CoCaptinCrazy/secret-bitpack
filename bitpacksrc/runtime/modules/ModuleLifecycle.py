from bitpack.runtime.modules.ModuleState import ModuleState


class ModuleLifecycle:

    def start(self, module):
        module.state.set(
            ModuleState.RUNNING
        )

    def stop(self, module):
        module.state.set(
            ModuleState.STOPPED
        )

    def crash(self, module):
        module.state.set(
            ModuleState.CRASHED
        )

    def restart(self, module):
        module.state.set(
            ModuleState.RUNNING
        )