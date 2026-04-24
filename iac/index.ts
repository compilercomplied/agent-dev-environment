import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { getAppConfig } from "./configuration";

const APP_ID = "agent-dev-environment";
const agentsNamespace = "agents";

const config = new pulumi.Config();
const appConfig = getAppConfig(config);

const configMap = new k8s.core.v1.ConfigMap("agent-dev-env-config", {
  metadata: {
    name: `${APP_ID}-configmap`,
    namespace: agentsNamespace
  },
  data: {
    ...appConfig.plainConfig,
  },
});

const secret = new k8s.core.v1.Secret("agent-dev-env-secret", {
  metadata: {
    name: `${APP_ID}-secret`,
    namespace: agentsNamespace,
  },
  stringData: {
    ...appConfig.secrets,
  },
});

const portStr = appConfig.plainConfig["AGENT_DEV_ENVIRONMENT_PORT"];
if (!portStr) {
  throw new Error("AGENT_DEV_ENVIRONMENT_PORT is required but not found in configuration");
}
const port = parseInt(portStr);
if (isNaN(port)) {
  throw new Error(`AGENT_DEV_ENVIRONMENT_PORT must be a valid number, got: ${portStr}`);
}

const service = new k8s.core.v1.Service(`${APP_ID}-svc`, {
  metadata: {
    name: APP_ID,
    namespace: agentsNamespace,
  },
  spec: {
    clusterIP: "None", // Headless service
    selector: { app: APP_ID }, // Pods must have this label to join the service
    ports: [{ port: port, targetPort: port }],
  },
});
