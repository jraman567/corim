#!/bin/bash
../cocli comid create --template SEV-Endorsement-Milan-Keys.json
../cocli comid create --template SEV-Endorsement-Referencevalues.json
../cocli corim create -m SEV-Endorsement-Milan-Keys.cbor -m SEV-Endorsement-Referencevalues.cbor -t corimMini.json -o sev-endorsement.cbor
