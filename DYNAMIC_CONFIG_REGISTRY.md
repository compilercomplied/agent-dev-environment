# Implementation Plan: Dynamic Agent Registry via Kubernetes ConfigMaps

## Objective
Enable `agent-hub` (Python) to dynamically discover the configuration schema (environment variables, defaults, and descriptions) of `agent-dev-environment` (Go). This allows AI agents to intelligently tweak deployment parameters based on self-documenting metadata.

## Architecture
1. **Source of Truth**: The `agent-dev-environment` Pulumi configuration (`Pulumi.<stack>.yaml`).
2. **Registry**: A Kubernetes ConfigMap named `agent-registry` in the `agents` namespace.
3. **Consumer**: The `K8sManager` in `agent-hub` fetches this ConfigMap to resolve templates before spawning pods.

---

## Step 1: Changes in `agent-dev-environment` (Go Project)

### 1.1 Define Structured Configuration
In the Go project's `Pulumi.local.yaml` (and other stacks), define the agent template as a structured object:

```yaml
config:
agent-dev-environment:
  image: ghcr.io/compilercomplied/agent-dev-environment:latest
  env_vars:
    AGENT_DEV_ENVIRONMENT_LOGGING_TYPE:
      default: "plain"
      description: "Format of the logs ('structured' or 'plain'). Plain is useful for local human-readable logs. Use structured for deployments."

### 1.2 Export to ConfigMap
In the Go project's `main.go` (Pulumi code), read this object and create the ConfigMap:

```go
conf := config.New(ctx, "agent-dev-environment")

type AgentTemplate struct {
    Image   string `json:"image"`
    EnvVars map[string]struct {
        Default     string `json:"default"`
        Description string `json:"description"`
    } `json:"env_vars"`
}

var templates map[string]AgentTemplate
conf.RequireObject("agent-templates", &templates)

_, err := corev1.NewConfigMap(ctx, "agent-registry", &corev1.ConfigMapArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("agent-registry"),
        Namespace: pulumi.String("agents"),
    },
    Data: pulumi.StringMap{
        "templates": pulumi.JSONMarshal(templates),
    },
})
```

---

## Step 2: Changes in `agent-hub` (Python Project)

### 2.1 Update `K8sManager` (`src/k8s/manager.py`)
Refactor the manager to fetch and apply these templates.

#### Add `get_template` Method:
- Use `self.core_v1.read_namespaced_config_map(name="agent-registry", namespace="agents")`.
- Parse the `templates` key from the `data` field (YAML/JSON).
- Return the dictionary for the requested `agent_name`.

#### Update `create_task` Method:
- **Signature**: `create_task(self, task: str, agent_name: str = "agent-dev-environment", overrides: dict[str, str] | None = None)`.
- **Merge Logic**:
    1. Fetch template.
    2. Iterate through `env_vars` in the template.
    3. Use value from `overrides` if present, otherwise use the `default` from the template.
    4. Construct `client.V1EnvVar` objects.

### 2.2 Update API Endpoint (`src/main.py`)
Update the `/api/v1/test-k8s-integration` (or create a new one) to allow passing configuration overrides. This is what the LLM will eventually call.

---

## Key Technical Constraints
- **Namespace**: Always use the `agents` namespace for the registry and the pods.
- **Image Pull Secrets**: Ensure `ghcr-secret-cddca01f` is included in the Pod Spec.
- **Validation**: Run `mise run static-analysis` after changes to ensure types and linting are correct.
- **Fail-Safe**: Provide hardcoded fallbacks in `get_template` in case the ConfigMap is missing during local development.
