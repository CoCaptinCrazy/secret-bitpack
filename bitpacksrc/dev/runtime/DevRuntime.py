import os


class DevRuntime:

    def launch(self, environment):

        print(
            f"Starting development environment: {environment.name}"
        )

        os.makedirs(
            f"data/dev/{environment.name}",
            exist_ok=True
        )

        print(
            f"Runtime: {environment.runtime}"
        )

        for service in environment.services:

            print(
                f"Starting service: {service}"
            )

        for port in environment.ports:

            print(
                f"Opening port: {port}"
            )

        for path in environment.watch:

            print(
                f"Watching: {path}"
            )

        return True
