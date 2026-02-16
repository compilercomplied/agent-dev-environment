import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";
import { getAppConfig } from "./libraries/configuration";

const APP_ID = "agent-dev-environment";
const agentsNamespace = "agents";

const config = new pulumi.Config();
const allConfig = pulumi.runtime.allConfig();
const appConfig = getAppConfig(config, "AGENT_DEV_ENVIRONMENT_", allConfig, config.name);

const configMap = new k8s.core.v1.ConfigMap("agent-dev-env-config", {
  metadata: {
    name: `${APP_ID}-configmap`,
    namespace: agentsNamespace
  },
  data: appConfig.plainConfig,
});

const secret = new k8s.core.v1.Secret("agent-dev-env-secret", {
  metadata: {
    name: `${APP_ID}-secret`,
    namespace: agentsNamespace,
  },
  stringData: {
      ...appConfig.secrets,
			// These are manually added; these values are not used by the go http api
			// but by the underlying binaries it executes.
      GITHUB_TOKEN: config.requireSecret("GITHUB_TOKEN"),
      PULUMI_CONFIG_PASSPHRASE: config.requireSecret("PULUMI_CONFIG_PASSPHRASE"),
  },
});
