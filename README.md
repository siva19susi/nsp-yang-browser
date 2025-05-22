# NSP YANG Browser

NSP allows you to create and execute intent-based automation model for implementing network-level changes that translates high-level designs into necessary network configuration. In order to faciliate this flow, YANG modeling is used to leverage the automation to its complete heights.

This tool provides a means to explore these YANG models either from a connected NSP system or the YANG repository can be manually uploaded and explored.

## Key Features

- **Path and Tree Browsing:** Easily navigate any YANG repository.
- **NSP Integration:** Directly browse YANG models of any connected NSP.
- **Sample Payloads:** View example payloads for each YANG model.
- **Model Comparison:** Seamlessly compare YANG models, whether manually uploaded or connected to an NSP.
- **Advanced Filtering:** Filter paths by prefix, defaults, and state attributes for more refined results.
- **Query NSP:** Query the connected NSP directly for live results of the yang paths.
- **Offline Mode:** View any YANG model in offline mode in a click of a button.

## Models that can be viewed

- NSP Modules
- LSO Operations
- Telemetry Types
- Intent Types

## How to try it out

```bash
mkdir ~/Downloads/offline

docker pull sivasusi19/nsp-yang-browser:latest

docker run -d -p 4173:4173 -v ~/Downloads/offline:/offline --name nsp-yang-browser sivasusi19/nsp-yang-browser:latest
```

You can access the Yang Browser from: `http://localhost:4173`
