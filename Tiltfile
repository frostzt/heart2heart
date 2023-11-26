#######################
# BUILD IMAGES
#######################

# Build keeper app image
docker_build(
  'keeper-service:local',
  dockerfile='./apps/keeper/Dockerfile.local',
  context='.',
  ignore=['./apps/bigboss/*', './apps/hippo/*', './apps/seer/*', './apps/summer/*', './apps/summer-e2e/*'],
  live_update=[sync('.', '/apps/keeper')]
)

#######################
# Resource Initialization
#######################

# Initialize Dapr
# local_resource('initialize dapr', cmd='dapr init -k', auto_init=True)

########################
# Deploy helm charts
#######################

# Deploy Alpha Service
k8s_yaml(helm('./infra/charts/keeper', name='keeper-local', values='./infra/charts/keeper/values.local.yaml'))
