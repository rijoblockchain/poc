const express = require("express");
const app = express();
var morgan = require('morgan')
app.use(morgan('combined'))
const bodyparser = require("body-parser");
require('dotenv').config();
const { registerUser, userExist } = require("./registerUser");
const {initiateGrant,assignGrant,acceptGrant,rejectGrant,revokeGrant,updateGrant,requestReimbursement,acceptReimbursement,rejectReimbursement,redeemTokens,acceptRedeem,rejectRedeem,addAwardee,addSubawardee,addProgress,deleteGrant} = require('./tx')
const {GetGrant,GetAllGrants,GetWallet,GetAllGrantsUser,GetAllApprovedGrants,GetGrantsByStatus,GetRemainingAmount,GetGrantBenefits,GetPayments,GetPaymentByAwardee,GetProgress,MyWallet,GetPaymentByStatus,GetPaymentByStatusForAllGrants,GetMSPIDs} =require('./query')
const PORT=process.env.PORT

var cors = require('cors')
app.use(cors())
app.use(bodyparser.json());

app.listen(PORT, () => {
    console.log(`server started at port ${PORT}`);

})

app.post("/register", async (req, res) => {

    try {
        let org = req.body.org[0].toUpperCase() + req.body.org.slice(1);
        let userId = req.body.userId;
        let result = await registerUser({ OrgMSP: org, userId: userId });
        res.send(result);

    } catch (error) {
        res.status(500).send(error)
    }
});


app.post("/initiateGrant", async (req, res) => {
    try {
        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "data": req.body.data
        }

        result = await initiateGrant(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/assignGrant", async (req, res) => {
    try {
        
        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "data": req.body.data
        }
        let result = await userExist({ OrgMSP: req.body.data.awardee[0].organization, userId: req.body.data.awardee[0].id });
        
        if (result != true) {
            res.send(result)
        } else {
            result = await assignGrant(payload);
            res.send(result)
        }

        
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/acceptGrant", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "id": req.body.id
        }

        let result = await acceptGrant(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/rejectGrant", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "id": req.body.id
        }

        let result = await rejectGrant(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/revokeGrant", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "id": req.body.id
        }

        let result = await revokeGrant(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/updateGrant", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "data": req.body.data
        }

        let result = await updateGrant(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/requestReimbursement", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "data": req.body.data
        }

        let result = await requestReimbursement(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/acceptReimbursement", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "grant_id": req.body.grant_id,
            "payment_id": req.body.payment_id
        }

        let result = await acceptReimbursement(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/rejectReimbursement", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "grant_id": req.body.grant_id,
            "payment_id": req.body.payment_id,
            "message": req.body.message
        }

        let result = await rejectReimbursement(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/redeemTokens", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "grant_id": req.body.grant_id,
            "payment_id": req.body.payment_id
        }

        let result = await redeemTokens(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/acceptRedeem", async (req, res) => {
    try {

        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "grant_id": req.body.grant_id,
            "payment_id": req.body.payment_id
        }

        let result = await acceptRedeem(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/rejectRedeem", async (req, res) => {
    try {

        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "grant_id": req.body.grant_id,
            "payment_id": req.body.payment_id,
            "message": req.body.message
        }

        let result = await rejectRedeem(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/addAwardee", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "data": req.body.data
        }

        let result = await userExist({ OrgMSP: req.body.data.awardee.organization, userId: req.body.data.awardee.id });
        
        if (result != true) {
            res.send(result)
        } else {
            result = await addAwardee(payload);
            res.send(result)
        }
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/addSubawardee", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "data": req.body.data
        }

        let result = await userExist({ OrgMSP: req.body.data.awardee.organization, userId: req.body.data.awardee.id });
        
        if (result != true) {
            res.send(result)
        } else {
            result = await addSubawardee(payload);
            res.send(result)
        }
    } catch (error) {
        res.status(500).send(error)
    }
})

app.post("/addProgress", async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "data": req.body.data
        }

        let result = await addProgress(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.get('/getWallet', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "grant_id": req.query.grantId,
            "awardee_id": req.query.awardeeId,
            "status": req.query.status
        }

        let result = await GetWallet(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/myWallet', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "grant_id": req.query.grantId,
        }

        let result = await MyWallet(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getGrant', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "id": req.query.id
        }

        let result = await GetGrant(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});



app.get('/getAllGrants', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId
        }

        let result = await GetAllGrants(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getAllGrantsUser', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId
        }

        let result = await GetAllGrantsUser(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getAllApprovedGrants', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId
        }

        let result = await GetAllApprovedGrants(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getGrantsByStatus', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "status": req.query.status
        }

        let result = await GetGrantsByStatus(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getPaymentByStatus', async (req, res) => {
   
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "status": req.body.status
        }

        let result = await GetPaymentByStatus(payload);
        res.json(result)
        
    } catch (error) {
        res.send(error)
    }
});

app.get('/getPaymentByAwardee', async (req, res) => {
   
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "grant_id": req.query.grantId,
            "awardee_id": req.query.awardeeId
        }

        let result = await GetPaymentByAwardee(payload);
        res.json(result)
        
    } catch (error) {
        res.send(error)
    }
});

app.get('/getPaymentByStatusForAllGrants', async (req, res) => {
    try {


        let payload = {
            "org": req.body.org[0].toUpperCase() + req.body.org.slice(1),
            "userId": req.body.userId,
            "status": req.body.status
        }

        let result = await GetPaymentByStatusForAllGrants(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getRemainingAmount', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "grant_id": req.query.grantId,
        }

        let result = await GetRemainingAmount(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getGrantBenefits', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "grantId": req.query.grantId,
        }

        let result = await GetGrantBenefits(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getPayments', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "grantId": req.query.grantId,
        }

        let result = await GetPayments(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.get('/getProgress', async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "grantId": req.query.grantId
        }

        let result = await GetProgress(payload);
        res.json(result)
    } catch (error) {
        res.send(error)
    }
});

app.post("/deleteGrant", async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId,
            "grantId": req.query.grantId
        }

        let result = await deleteGrant(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})

app.get("/getMSPIDs", async (req, res) => {
    try {


        let payload = {
            "org": req.query.org[0].toUpperCase() + req.query.org.slice(1),
            "userId": req.query.userId
        }

        let result = await GetMSPIDs(payload);
        res.send(result)
    } catch (error) {
        res.status(500).send(error)
    }
})
