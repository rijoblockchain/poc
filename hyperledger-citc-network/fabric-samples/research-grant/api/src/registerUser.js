const { Wallets } = require("fabric-network");
const FabricCAServices = require('fabric-ca-client');

const { buildCAClient, registerAndEnrollUser, enrollAdmin ,userExist} = require("./CAUtil")
const {  buildWallet } = require("./AppUtils");
const { getCCP } = require("./buildCCP");
const path=require('path');
const { Utils: utils } = require('fabric-common');
let config=utils.getConfig()
config.file(path.resolve(__dirname,'config.json'))
let walletPath;
exports.registerUser = async ({ OrgMSP, userId }) => {

    let org = OrgMSP.replace('MSP','').toLowerCase();;
    walletPath=path.join(__dirname,`wallet/${OrgMSP}`)
    let ccp = getCCP(OrgMSP)
    const caClient = buildCAClient(FabricCAServices, ccp, `ca.${org}.example.com`);

    // setup the wallet to hold the credentials of the application user
    const wallet = await buildWallet(Wallets, walletPath);

    // in a real application this would be done on an administrative flow, and only once
    await enrollAdmin(caClient, wallet, OrgMSP);

    // in a real application this would be done only when a new user was required to be added
    // and would be part of an administrative flow
    let result = await registerAndEnrollUser(caClient, wallet, OrgMSP, userId, `${org}.department1`);

    return result;
}


exports.userExist=async({ OrgMSP, userId })=>{
    let org = OrgMSP.replace('MSP','').toLowerCase();
    OrgMSP = OrgMSP.replace('MSP','')
    walletPath=path.join(__dirname,`wallet/${OrgMSP}`)
    let ccp = getCCP(OrgMSP)
    const caClient = buildCAClient(FabricCAServices, ccp, `ca.${org}.example.com`);

    // setup the wallet to hold the credentials of the application user
    const wallet = await buildWallet(Wallets, walletPath);

   const result=await userExist(wallet,userId)
   return result;
}
