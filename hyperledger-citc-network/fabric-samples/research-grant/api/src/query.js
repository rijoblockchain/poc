const { getCCP } = require("./buildCCP");
const { Wallets, Gateway } = require('fabric-network');
const path = require("path");
const { buildWallet } = require('./AppUtils');
const { IdentityService } = require("fabric-ca-client");

const chaincodeName = process.env.chaincodeName;
const channelName = process.env.channelName

exports.GetWallet = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);
    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);
    let result = await contract.evaluateTransaction("GetWallet", request.grant_id, request.awardee_id, request.status);
    return JSON.parse(result);
}

exports.MyWallet = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);
    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);
    let result = await contract.evaluateTransaction("MyWallet", request.grant_id);
    return JSON.parse(result);
}

exports.GetGrant = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);
    let id = request.id;
    let result = await contract.evaluateTransaction("ReadGrant", id);
    return JSON.parse(result);
}

exports.GetAllGrants = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);

    let result = await contract.evaluateTransaction("GetAllGrants");
    return JSON.parse(result);
}

exports.GetAllGrantsUser = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);

    let result = await contract.evaluateTransaction("GetAllGrantsUser");
    return JSON.parse(result);
}

exports.GetAllApprovedGrants = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);

    let result = await contract.evaluateTransaction("GetAllApprovedGrants");
    return JSON.parse(result);
}

exports.GetGrantsByStatus = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);

    let result = await contract.evaluateTransaction("GetGrantsByStatus", request.status);
    return JSON.parse(result);
}

exports.GetPaymentByStatusForAllGrants = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);

    let result = await contract.evaluateTransaction("GetPaymentByStatusForAllGrants", request.status);
    return JSON.parse(result);
}


exports.GetRemainingAmount = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);
    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);
    let result = await contract.evaluateTransaction("GetRemainingAmount", request.grant_id);
    return JSON.parse(result);
}


exports.GetGrantBenefits = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);

    let result = await contract.evaluateTransaction("GetGrantBenefits", request.grantId);
    return JSON.parse(result);
}

exports.GetPayments = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);

    let result = await contract.evaluateTransaction("GetPayments", request.grantId);
    return JSON.parse(result);
}

exports.GetPaymentByStatus = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);
    let status = JSON.stringify(request.status)
    let result = await contract.evaluateTransaction("GetPaymentByStatus", status);
    return JSON.parse(result);
    
}

exports.GetPaymentByAwardee = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);
    let result = await contract.evaluateTransaction("GetPaymentByAwardee", request.grant_id, request.awardee_id);
    return JSON.parse(result);
    
}

exports.GetPaymentByStatusForAllGrants = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);
    let status = JSON.stringify(request.status)
    let result = await contract.evaluateTransaction("GetPaymentByStatusForAllGrants", status);
    return JSON.parse(result);
    
}

exports.GetProgress = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    // Get the contract from the network.
    const contract = network.getContract(chaincodeName);

    let result = await contract.evaluateTransaction("GetProgress", request.grantId);
    return JSON.parse(result);
}

exports.GetMSPIDs = async (request) => {
    let org = request.org;
    const walletPath = path.join(__dirname,`wallet/${org}`)
    const ccp = getCCP(org);

    const wallet = await buildWallet(Wallets, walletPath);

    const gateway = new Gateway();

    await gateway.connect(ccp, {
        wallet,
        identity: request.userId,
        discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
    });

    // Build a network instance based on the channel where the smart contract is deployed
    const network = await gateway.getNetwork(channelName);

    const channel = network.getChannel()
    const result = channel.getMspids()
    return result;
}