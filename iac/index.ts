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
