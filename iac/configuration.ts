import * as pulumi from "@pulumi/pulumi";

export interface AppConfig {
  secrets: Record<string, pulumi.Input<string>>;
  plainConfig: Record<string, string>;
}

export function getAppConfig(config: pulumi.Config): AppConfig {
  const allCfg = pulumi.runtime.allConfig();
  // PULUMI_CONFIG_SECRET_KEYS is set by the Pulumi CLI (and Automation API via setAllConfig)
  // and lists the fully-qualified keys (project:key) that contain secret values.
  const secretKeySet = new Set<string>(
    JSON.parse(process.env["PULUMI_CONFIG_SECRET_KEYS"] ?? "[]")
  );

  const secrets: Record<string, pulumi.Input<string>> = {};
  const plainConfig: Record<string, string> = {};
  const prefix = `${config.name}:`;

  for (const [key, value] of Object.entries(allCfg)) {
    if (!key.startsWith(prefix)) continue;

    const varName = key.slice(prefix.length);
    if (secretKeySet.has(key)) {
      secrets[varName] = pulumi.secret(value);
    } else {
      plainConfig[varName] = value;
    }
  }

  return { secrets, plainConfig };
}
