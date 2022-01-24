# rc-contract-go

This is the smart contract based on golang.

The smart contract can be tested with the following steps.

### Bring up the test network

Open a terminal from the /test-network dictionary, run

```
./network.sh down
./network.sh up createChannel
``` 

Add the binaries and config filepath to the CLI path

```
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
```

