const pulumi = require("@pulumi/pulumi");
const allConfig = pulumi.runtime.allConfig();
console.log(allConfig);
