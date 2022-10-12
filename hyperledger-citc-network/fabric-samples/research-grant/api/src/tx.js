const { getCCP } = require("./buildCCP");
const { Wallets, Gateway } = require('fabric-network');
const path = require("path");
const {buildWallet} =require('./AppUtils');

const chaincodeName = process.env.chaincodeName;
const channelName = process.env.channelName

let gateway
/*
{
    org:Org1MSP,
    channelName:"mychannel",
    chaincodeName:"basic",
    userId:"aditya"
    data:{
        id:"asset1",
        color:"red",
        size:5,
        appraisedValue:200,
        owner:"TOM"
    }
}

*/
exports.initiateGrant = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let statefulTxn = contract.createTransaction('InitiateGrant');
            let data=request.data;
            let tmapData = Buffer.from(JSON.stringify(data));
            statefulTxn.setTransient({
                grant: tmapData
            });
            let result = await statefulTxn.submit();
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }
    
   
}

exports.assignGrant = async (request) => {
    try{
        
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let statefulTxn = contract.createTransaction('AssignGrant');
            let data=request.data;
            let tmapData = Buffer.from(JSON.stringify(data));
            statefulTxn.setTransient({
                assign_grant: tmapData
            });
            let result = await statefulTxn.submit();
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }
}

exports.acceptGrant = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let id=request.id;
            let result = await contract.submitTransaction('AcceptGrant',id);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }
}

exports.rejectGrant = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let id=request.id;
            let result = await contract.submitTransaction('RejectGrant',id);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }
    
   
}

exports.revokeGrant = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let id=request.id;
            grantorId=request.userId;
            let result = await contract.submitTransaction('RevokeGrant',id, grantorId);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }
    
   
}

exports.updateGrant = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let statefulTxn = contract.createTransaction('UpdateGrant');
            let data=request.data;
            let tmapData = Buffer.from(JSON.stringify(data));
            statefulTxn.setTransient({
                update_grant: tmapData
            });
            let result = await statefulTxn.submit();
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }  
}

exports.requestReimbursement = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let statefulTxn = contract.createTransaction('RequestReimbursement');
            let data=request.data;
            let tmapData = Buffer.from(JSON.stringify(data));
            statefulTxn.setTransient({
                request_reimbursement: tmapData
            });
            let result = await statefulTxn.submit();
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }  
}

exports.acceptReimbursement = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let grant_id=request.grant_id;
            let payment_id=request.payment_id;
            let result = await contract.submitTransaction('AcceptReimbursement',grant_id, payment_id);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }   
}

exports.rejectReimbursement = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let grant_id=request.grant_id;
            let payment_id=request.payment_id;
            let message=request.message;
            let result = await contract.submitTransaction('RejectReimbursement',grant_id, payment_id,message);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }   
}

exports.redeemTokens = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let grant_id=request.grant_id;
            let payment_id=request.payment_id;
            let result = await contract.submitTransaction('RedeemTokens',grant_id, payment_id);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }   
}

exports.acceptRedeem = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let grant_id=request.grant_id;
            let payment_id=request.payment_id;
            let result = await contract.submitTransaction('AcceptRedeem',grant_id, payment_id);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }   
}

exports.rejectRedeem = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let grant_id=request.grant_id;
            let payment_id=request.payment_id;
            let message=request.message;
            let result = await contract.submitTransaction('rejectRedeem',grant_id, payment_id, message);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }   
}

exports.addSubawardee = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let statefulTxn = contract.createTransaction('AddSubawardee');
            let data=request.data;
            let tmapData = Buffer.from(JSON.stringify(data));
            statefulTxn.setTransient({
                add_subawardee: tmapData
            });
            let result = await statefulTxn.submit();
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }
    
   
}

exports.addAwardee = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let statefulTxn = contract.createTransaction('AddAwardee');
            let data=request.data;
            let tmapData = Buffer.from(JSON.stringify(data));
            statefulTxn.setTransient({
                add_awardee: tmapData
            });
            let result = await statefulTxn.submit();
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }   
}

exports.addProgress = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let statefulTxn = contract.createTransaction('AddProgress');
            let data=request.data;
            let tmapData = Buffer.from(JSON.stringify(data));
            statefulTxn.setTransient({
                add_progress: tmapData
            });
            let result = await statefulTxn.submit();
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }   
}

exports.deleteGrant = async (request) => {
    try{
        let org = request.org;
        const walletPath = path.join(__dirname,`wallet/${org}`)
        const ccp = getCCP(org);
    
        const wallet = await buildWallet(Wallets, walletPath);
    
        gateway = new Gateway();
    
        await gateway.connect(ccp, {
            wallet,
            identity: request.userId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
    
        // Build a network instance based on the channel where the smart contract is deployed
        const network = await gateway.getNetwork(channelName);
    
        // Get the contract from the network.
        const contract = network.getContract(chaincodeName);
    
        try {
            let transaction = contract.createTransaction('DeleteGrant');
			let result =  await transaction.submit(request.grantId);
            const response = {
                status: result.toString()
            }
            return (response);
    
        } catch (error) {
            console.log(`   Successfully caught the error: \n    ${error}`);
            const response = {
                status: 'error',
                message: error.message.split('message=').pop()
            }
            return (response)
            
        } 
    } finally {
        // Disconnect from the gateway peer when all work for this client identity is complete
        gateway.disconnect();
    }
    
   
}
