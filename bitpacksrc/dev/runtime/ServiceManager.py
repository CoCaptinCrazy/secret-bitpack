class ServiceManager:

    def start(self, service):

        print(
            f"[SERVICE] Starting {service}"
        )

        return True

    def stop(self, service):

        print(
            f"[SERVICE] Stopping {service}"
        )

        return True

    def restart(self, service):

        self.stop(service)

        self.start(service)

    def start_all(self, services):

        for service in services:

            self.start(service)
