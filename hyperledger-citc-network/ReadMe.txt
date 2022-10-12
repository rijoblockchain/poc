This setup was created on Ubuntu 20.04. It uses the Fablo tool (https://github.com/hyperledger-labs/fablo) to quickly create a Hyperledger network in Docker environment. Please perform the following steps:
1- Run the script ./hyperledger-citc.sh to install some necessary software 
2- Run the command sudo ./fablo recreate to start the network
3- You may open the file  fablo-config.json to view the network components. It uses the cryptogen tools to create crypto files
4- You may use the commands sudo ./fablo [down | start | stop | up | prune | reset] to interact with the network
5- The connection profile files are created in hyperledger-citc-network/fablo-target/fabric-config/connection-profile
6- When using the api to create users, it fails because the connection profile uses "url": "http://localhost:7040", so I changed it to "https"
7- Now, it fails:
"Failed to register user : Error: Calling register endpoint failed with error [Error: unable to verify the first certificate]"
8- How can we fix the api that it works with cryptoget certs?




