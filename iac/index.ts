import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { getAppConfig } from "./libraries/configuration";

const config = new pulumi.Config();
const allConfig = pulumi.runtime.allConfig();
const appConfig = getAppConfig(config, "AGENT_DEV_ENVIRONMENT_", allConfig, config.name);

const configMap = new k8s.core.v1.ConfigMap("agent-dev-env-config", {
  metadata: {
    name: "agent-dev-env-config",
  },
  data: appConfig.plainConfig,
});

const secret = new k8s.core.v1.Secret("agent-dev-env-secret", {
  metadata: {
    name: "agent-dev-env-secret",
  },
  stringData: appConfig.secrets,
});
