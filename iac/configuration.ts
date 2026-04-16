import * as pulumi from "@pulumi/pulumi";

export interface EnvVarDefinition {
  default: any;
  description: string;
}

export interface AppConfig {
  secrets: { [key: string]: pulumi.Output<string> };
  plainConfig: { [key: string]: string };
  envVars: Record<string, any>;
}

export function getAppConfig(config: pulumi.Config): AppConfig {
  const envVars = config.requireObject<Record<string, EnvVarDefinition>>("env_vars");
  const secrets: Record<string, pulumi.Output<string>> = {};
  const plainConfig: Record<string, string> = {};
  const processedEnvVars: Record<string, any> = {};

  const allConfig = pulumi.runtime.allConfig();
  const raw = JSON.parse(allConfig[`${config.name}:env_vars`]);

  for (const key of Object.keys(envVars)) {
    if (raw[key]?.default?.secure) {
      const secretValue = config.requireSecretObject<any>("env_vars").apply(ev => String(ev[key].default));
      secrets[key] = secretValue;
      processedEnvVars[key] = {
        ...envVars[key],
        default: secretValue,
      };
    } else {
      plainConfig[key] = String(envVars[key].default);
      processedEnvVars[key] = envVars[key];
    }
  }

  return { secrets, plainConfig, envVars: processedEnvVars };
}
