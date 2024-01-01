#######################
# BUILD IMAGES
#######################

# Build keeper app image
docker_build(
  'keeper-service:local',
  dockerfile='./apps/keeper/Dockerfile.local',
  context='.',
  live_update=[sync('.', '/apps/keeper')]
)

#######################
# Resource Initialization
#######################

# Initialize Dapr
# local_resource('initialize dapr', cmd='dapr init -k', auto_init=True)

local_resource('initialize signoz...', cmd='helm --namespace platform install signoz-apm signoz/signoz', auto_init=True)

########################
# Deploy helm charts
#######################

# Deploy Keeper Service
k8s_yaml(helm('./infra/charts/keeper', name='keeper-local', values='./infra/charts/keeper/values.local.yaml'))
k8s_resource('keeper-service', new_name='keeper-resources', labels=['keeper'])

k8s_yaml(helm('./infra/charts/signoz', name='signoz-local', values='./infra/charts/signoz/values.local.yaml'))
