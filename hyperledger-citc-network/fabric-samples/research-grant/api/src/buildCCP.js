const { buildCCPGrantor,buildCCPAwardee,buildCCPAuditor,buildCCPSubawardee} = require("./AppUtils");

exports.getCCP = (org) => {
    let ccp;
    switch (org) {
        case "Grantor":
            ccp = buildCCPGrantor();
            break;
        case "Awardee":
            ccp = buildCCPAwardee();
            break;
        case "Auditor":
            ccp = buildCCPAuditor();
            break;
        case "Subawardee":
            ccp = buildCCPSubawardee();
            break;
    }
    return ccp;
}