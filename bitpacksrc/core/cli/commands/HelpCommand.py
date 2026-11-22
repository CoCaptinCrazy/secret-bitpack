class HelpCommand:

    def execute(self, args):

        print("""

GENERAL
------------------------------------------------
  help             Show help
  version          Show version information


BLUEPRINTS
------------------------------------------------
  forge            Build a blueprint
  verify           Verify a blueprint
  blueprints       List blueprints
  inspect          Inspect a blueprint
  export           Export a blueprint
  import           Import a blueprint


RUNTIME
------------------------------------------------
  ignite           Run a blueprint
  modules          List running modules
  stop             Stop a module
  restart          Restart a module
  logs             View module logs


SNAPSHOTS
------------------------------------------------
  snapshot         Create a snapshot
  restore          Restore a snapshot


ASSEMBLIES
------------------------------------------------
  assemble         Deploy a .bpack assembly
  assemblies       List assemblies
  destroy          Destroy an assembly
  scale            Scale a module
  rollout          Roll out an update
  rollback         Roll back an update


REGISTRY
------------------------------------------------
  publish          Upload blueprint
  acquire          Download blueprint
  search           Search registry
  tags             List tags


AUTHENTICATION
------------------------------------------------
  login            Login to registry
  logout           Logout from registry


DEVELOPMENT
------------------------------------------------
  dev              Launch a .bpdev environment
  shell            Open development shell
  watch            Watch files for changes


BATCH
------------------------------------------------
  batch            Execute a .bpbatch workflow
  batch-status     Show batch status


NETWORKING
------------------------------------------------
  mesh             Manage a mesh
  meshes           List meshes
  bridges          List bridges
  route            Manage routes


SECURITY
------------------------------------------------
  health           Health checks
  policy           Manage policies
  alert            Create an alert
  secrets          Manage secrets


MONITORING
------------------------------------------------
  metrics          View metrics
  profile          Profile runtime
  alerts           List alerts
  events           View events


CLUSTERING
------------------------------------------------
  cluster          Manage cluster
  clusters         List clusters
  nodes            List nodes


SYSTEM
------------------------------------------------
  runtime          Runtime information
  config           Manage configuration
  backup           Create backups


PIPELINES
------------------------------------------------
  pipeline         Execute a pipeline


FILE TYPES
------------------------------------------------
  Bitfile          Blueprint source definition
  .bitignore       Forge ignore rules
  .bp              Blueprint artifact
  .bpdev           Development environment
  .bpack           Assembly topology
  .bpbatch         Workflow batch file

===============================================
""")