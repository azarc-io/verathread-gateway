# load extensions
load('ext://dotenv', 'dotenv')
load('ext://helm_resource', 'helm_resource')
load('ext://dotenv', 'dotenv')
load('ext://restart_process', 'docker_build_with_restart')
load('ext://helm_resource', 'helm_resource')

# .env support
dotenv(fn='.env')

# tilt_config.json support
config.define_bool("hmr")
config.define_string("arch")
config.define_string_list("to-debug")
config.define_string("gateway-debug-port")
config.define_string("kube_context")
cfg = config.parse()

# arch type to build for, default is amd64, you can update this in the tilt_config.json file
arch = cfg.get('arch', 'amd64')

# debug ports
gatewayDebugPort = cfg.get("gateway-debug-port", "40000")

# enabling hmr will cause the gateway to stop proxying the front end directly and instead the ingress
# will point to the front end on the existing host, the user will notice no difference but hot reloading
# will be enabled and the build will run in kubernetes instead of your local machine
hmr = cfg.get('hmr', False)

# what kube context to permit, prevents you from switching your kube context to staging for eg. and mistakenly deploying there
allow_k8s_contexts(cfg.get('kube_context', 'k3d-dev-1'))

# default registry configuration
default_registry(
    'k3d-local-registry:5000',
    host_from_cluster='localhost:5000'
)

local_resource(
  'gateway-compile',
  'CGO_ENABLED=0 GOOS=linux GOARCH=%s go build -gcflags "all=-N -l" -o ./bin/vth-gateway ./cmd/app/main.go' % arch,
  deps=['./cmd/app/main.go', './cmd/app', './internal/', './pkg/'],
  labels=["compile"],
  resource_deps=[]
)

# enabling debug for an app will consume more memory because dlv will be used to enable remote debugging
# so enable selectively in the tilt_config.json file
to_debug = cfg.get('to-debug', [])

# the default entry path
entrypoint = '/vth-gateway'
appDockerFile = './deployment/docker/app/Dockerfile.tilt'
webDockerFile = './deployment/docker/web/Dockerfile.tilt'

# entry path to use if debug is enabled
if 'gateway' in to_debug:
    entrypoint = '/dlv --listen=:40000 --api-version=2 --headless=true --only-same-user=false --accept-multiclient exec --continue /vth-gateway'

# dockerfile to use if hmr is enabled
if hmr:
    webDockerFile = './deployment/docker/web/Dockerfile.hmr.tilt'

# watches directories for changes and triggers an update of the docker image
docker_build_with_restart(
  'vth-gateway',
  context='.',
  entrypoint=entrypoint,
  dockerfile=appDockerFile,
  platform='linux/%s' % arch,
  only=[
    './bin',
  ],
  live_update=[
    sync('./bin/', '/'),
  ]
)

# hmr mode will sync the source files to the sidecar and the build will happen there, in this mode
# the gateway will not proxy but instead the ingress will point at the side car for / where as anything
# else such as /graphql or /api would point at the gateway
if hmr:
    docker_build(
      'vth-gateway-web',
      context='.',
      entrypoint='vite build --watch --outDir /web',
      dockerfile=webDockerFile,
      platform='linux/%s' % arch,
      only=[
        './cmd/web',
      ],
      ignore=[
        './cmd/web/node_modules',
        './cmd/web/dist'
      ],
      live_update=[
        fall_back_on(['./cmd/web2/package.json', './cmd/web2/yarn.lock']),
        sync('./cmd/web', '/src'),
      ]
    )
    # Forward hmr port
    k8s_resource(workload='gateway-chart', port_forwards=3000)
# non hmr mode will simply sync your local dist folder, you can use vite build --watch
# in this mode the gateway will serve the files and the side car becomes an ephemeral initContainer
else:
    docker_build_with_restart(
      'vth-gateway-web',
      context='.',
      entrypoint='echo "reloading" && cp -R /static/. /web/ && tail -f /dev/null',
      dockerfile=webDockerFile,
      platform='linux/%s' % arch,
      only=[
        './cmd/web/dist',
      ],
      live_update=[
        sync('./cmd/web/dist', '/static'),
      ]
    )

# watches the chart directory and triggers an update if yaml files change
helm_resource(
  'gateway-chart',
  './deployment/charts/gateway',
  namespace=os.getenv('NAMESPACE'),
  deps=["./deployment/charts/gateway"],
  flags=[
    '--set=configuration.auth.client_id=%s' % os.getenv('AUTH_CLIENT_ID'),
    '--set=configuration.auth.client_secret=%s' % os.getenv('AUTH_CLIENT_SECRET'),
    '--set=configuration.auth.domain=%s' % os.getenv('AUTH_DOMAIN'),
    '--set=dev=true',
    '--set=bind.debug=%s'% gatewayDebugPort,
    '--set=hmr=%s' % hmr
  ],
  image_deps=['vth-gateway', 'vth-gateway-web'],
  image_keys=[('image.repository', 'image.tag'), ('image.web_repository', 'image.web_repository_tag')],
)

# forward gateway public port
k8s_resource(workload='gateway-chart', port_forwards=6010)

# Forward mongo port
#local_resource(
#  name='mongo',
#  serve_cmd='kubectl port-forward service/mongodb 27017:27017 --namespace=%s' % os.getenv('NAMESPACE'),
#  labels=['ports'],
#  links=['mongo://localhost:27017'],
#  readiness_probe=probe(tcp_socket=tcp_socket_action(port=27017))
#)

if 'gateway' in to_debug:
    k8s_resource('gateway-chart',
        port_forwards=[
            gatewayDebugPort,  # debugger
        ],
        labels=['deployment'],
        links=[
            'localhost:%s' % gatewayDebugPort,
        ]
    )
