#!/bin/bash
colorcli tx send cosmos1wqv8jxzseyvmk93gpa6qxdm9l9cwpj3htkey2z $1 100stake --chain-id=testing --generate-only > unsignedSendTx.json
echo 12345678 |colorcli tx sign   --chain-id=testing --from=cosmos1wqv8jxzseyvmk93gpa6qxdm9l9cwpj3htkey2z  unsignedSendTx.json > signedSendTx.json
colorcli tx broadcast signedSendTx.json
