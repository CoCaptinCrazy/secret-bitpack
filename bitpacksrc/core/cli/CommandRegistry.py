from bitpack.core.cli.commands.HelpCommand import HelpCommand
from bitpack.core.cli.commands.VersionCommand import VersionCommand

from bitpack.core.cli.commands.ForgeCommand import ForgeCommand
from bitpack.core.cli.commands.VerifyCommand import VerifyCommand
from bitpack.core.cli.commands.BlueprintsCommand import BlueprintsCommand
from bitpack.core.cli.commands.InspectCommand import InspectCommand
from bitpack.core.cli.commands.ExportCommand import ExportCommand
from bitpack.core.cli.commands.ImportCommand import ImportCommand

from bitpack.core.cli.commands.IgniteCommand import IgniteCommand
from bitpack.core.cli.commands.ModulesCommand import ModulesCommand
from bitpack.core.cli.commands.StopCommand import StopCommand
from bitpack.core.cli.commands.RestartCommand import RestartCommand
from bitpack.core.cli.commands.LogsCommand import LogsCommand

from bitpack.core.cli.commands.SnapshotCommand import SnapshotCommand
from bitpack.core.cli.commands.RestoreCommand import RestoreCommand

from bitpack.core.cli.commands.AssembleCommand import AssembleCommand
from bitpack.core.cli.commands.AssembliesCommand import AssembliesCommand
from bitpack.core.cli.commands.DestroyCommand import DestroyCommand
from bitpack.core.cli.commands.ScaleCommand import ScaleCommand
from bitpack.core.cli.commands.RolloutCommand import RolloutCommand
from bitpack.core.cli.commands.RollbackCommand import RollbackCommand

from bitpack.core.cli.commands.PublishCommand import PublishCommand
from bitpack.core.cli.commands.AcquireCommand import AcquireCommand
from bitpack.core.cli.commands.SearchCommand import SearchCommand
from bitpack.core.cli.commands.TagsCommand import TagsCommand

from bitpack.core.cli.commands.LoginCommand import LoginCommand
from bitpack.core.cli.commands.LogoutCommand import LogoutCommand

from bitpack.dev.command.DevCommand import DevCommand
from bitpack.core.cli.commands.ShellCommand import ShellCommand
from bitpack.core.cli.commands.WatchCommand import WatchCommand

from bitpack.batch.command.BatchCommand import BatchCommand
from bitpack.core.cli.commands.BatchStatusCommand import BatchStatusCommand

from bitpack.core.cli.commands.MeshCommand import MeshCommand
from bitpack.core.cli.commands.MeshesCommand import MeshesCommand
from bitpack.core.cli.commands.BridgesCommand import BridgesCommand
from bitpack.core.cli.commands.RouteCommand import RouteCommand

from bitpack.core.cli.commands.HealthCommand import HealthCommand
from bitpack.core.cli.commands.PolicyCommand import PolicyCommand
from bitpack.core.cli.commands.AlertCommand import AlertCommand

from bitpack.core.cli.commands.MetricsCommand import MetricsCommand
from bitpack.core.cli.commands.ProfileCommand import ProfileCommand
from bitpack.core.cli.commands.AlertsCommand import AlertsCommand
from bitpack.core.cli.commands.EventsCommand import EventsCommand

from bitpack.core.cli.commands.SecretsCommand import SecretsCommand

from bitpack.core.cli.commands.ClusterCommand import ClusterCommand
from bitpack.core.cli.commands.ClustersCommand import ClustersCommand
from bitpack.core.cli.commands.NodesCommand import NodesCommand

from bitpack.core.cli.commands.RuntimeCommand import RuntimeCommand
from bitpack.core.cli.commands.ConfigCommand import ConfigCommand
from bitpack.core.cli.commands.BackupCommand import BackupCommand

from bitpack.core.cli.commands.PipelineCommand import PipelineCommand


class CommandRegistry:

    def __init__(self):

        self.commands = {

            # General
            "help": HelpCommand(),
            "version": VersionCommand(),

            # Blueprints
            "forge": ForgeCommand(),
            "verify": VerifyCommand(),
            "blueprints": BlueprintsCommand(),
            "inspect": InspectCommand(),
            "export": ExportCommand(),
            "import": ImportCommand(),

            # Runtime
            "ignite": IgniteCommand(),
            "modules": ModulesCommand(),
            "stop": StopCommand(),
            "restart": RestartCommand(),
            "logs": LogsCommand(),

            # Snapshots
            "snapshot": SnapshotCommand(),
            "restore": RestoreCommand(),

            # Assemblies
            "assemble": AssembleCommand(),
            "assemblies": AssembliesCommand(),
            "destroy": DestroyCommand(),
            "scale": ScaleCommand(),
            "rollout": RolloutCommand(),
            "rollback": RollbackCommand(),

            # Registry
            "publish": PublishCommand(),
            "acquire": AcquireCommand(),
            "search": SearchCommand(),
            "tags": TagsCommand(),

            # Authentication
            "login": LoginCommand(),
            "logout": LogoutCommand(),

            # Development
            "dev": DevCommand(),
            "shell": ShellCommand(),
            "watch": WatchCommand(),

            # Batch
            "batch": BatchCommand(),
            "batch-status": BatchStatusCommand(),

            # Networking
            "mesh": MeshCommand(),
            "meshes": MeshesCommand(),
            "bridges": BridgesCommand(),
            "route": RouteCommand(),

            # Security
            "health": HealthCommand(),
            "policy": PolicyCommand(),
            "alert": AlertCommand(),
            "secrets": SecretsCommand(),

            # Monitoring
            "metrics": MetricsCommand(),
            "profile": ProfileCommand(),
            "alerts": AlertsCommand(),
            "events": EventsCommand(),

            # Clustering
            "cluster": ClusterCommand(),
            "clusters": ClustersCommand(),
            "nodes": NodesCommand(),

            # System
            "runtime": RuntimeCommand(),
            "config": ConfigCommand(),
            "backup": BackupCommand(),

            # Pipelines
            "pipeline": PipelineCommand()
        }

    def register(
        self,
        name,
        command
    ):
        self.commands[name] = command

    def unregister(
        self,
        name
    ):
        if name in self.commands:
            del self.commands[name]

    def get(
        self,
        command
    ):
        return self.commands.get(command)

    def exists(
        self,
        command
    ):
        return command in self.commands

    def list(self):
        return sorted(
            self.commands.keys()
        )