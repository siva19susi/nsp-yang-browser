# YANG Browser

Browser to simplify YANG model visualization, eliminating the hassle of juggling multiple tools to access this information.

## Key Features

- **Path and Tree Browsing:** Easily navigate any YANG repository.
- **NSP Integration:** Directly browse YANG models of any connected NSP.
- **Sample Payloads:** View example payloads for each YANG model.
- **Model Comparison:** Seamlessly compare YANG models, whether manually uploaded or connected to an NSP.
- **Advanced Filtering:** Filter paths by prefix, defaults, and state attributes for more refined results.

## How to try it out

```bash
mkdir ~/Downloads/uploads

docker pull sivasusi19/yang-browser:latest

docker run -d -p 4173:4173 -v ~/Downloads/uploads:/uploads --name yang-browser sivasusi19/yang-browser:latest
```

You can access the Yang Browser from: `http://localhost:4173`
