from bitpack.dev.parser.DevEnvironment import DevEnvironment


class DevParser:

    def parse(self, source):

        env = DevEnvironment()

        lines = source.splitlines()

        for raw in lines:

            line = raw.strip()

            if not line:
                continue

            if line.startswith("development"):

                parts = line.split('"')

                if len(parts) >= 2:
                    env.name = parts[1]

            elif line.startswith("using "):

                env.runtime = line.split()[1]

            elif line.startswith("workspace "):

                env.workspace = (
                    line
                    .replace("workspace", "")
                    .replace('"', "")
                    .strip()
                )

            elif line.startswith("service "):

                env.services.append(
                    line.split()[1]
                )

            elif line.startswith("port "):

                env.ports.append(
                    int(line.split()[1])
                )

            elif line.startswith("watch "):

                env.watch.append(
                    line.split()[1]
                )

        return env