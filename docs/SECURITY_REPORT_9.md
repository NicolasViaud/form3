# Activate TLS for REST API

***
<font color="#df3d03">Severity : 8.1 High</font>  
Date : 04/06/2023  
Reporter : Nicolas Viaud   
Weakness enumeration : [CWE-319: Cleartext Transmission of Sensitive Information](https://cwe.mitre.org/data/definitions/319.html)  
CVSS v3.1 Vector : `AV:A/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:N`
***

According to the README.md, the configuration in the `k8s/` directory is used in production:
> k8s/ folder contains all the necessary k8s manifests to run this application in a production environment.

However, in the current configuration, the communications between a requester and the API servers are not secure because they are using the protocol `http` instead of `https`. There is 2 major flaws due to this insecure communication:
- data between the requester and the API server are not encrypted. An attacker could see all the data transiting over the network. Some data, like the JWT token are considered sensitive data
- the identity of the server could be hijacked by an attacker pretending to be the API servers. He can intercept the requests and modify them before sending them back to the legitimate server or just read the data of the request

To create a secure communication channel between a requester and the innsecure API, the server need to use a secure version of the TLS protocol. This TLS protocol will:
- authenticate the server, based on the trust of a third party authority
- encrypted the communication
- verify the integrity of the data after an exchange

## Developer Fix

A [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) has to be configured. It acts as a reverse proxy to route the request to the innsecure service API and decrypt the data sent over https.
The configuration of the Ingress will depend on the hosting platform. Because the hosting platform is not known, we can assume to configure the ingress throw a nginx server. The documentation can be found here: https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/

### Security requirements

The list of TLS version supported are:
- V1.3 (default)
- V1.2

ALL the other TLS version are deprecated and shouldn't be supported by the ingress.

The list of supported cypher to configure for the TLS V1.3 are :
- TLS_AES_256_GCM_SHA384
- TLS_AES_128_CCM_SHA256
- TLS_CHACHA20_POLY1305_SHA256

The list of supported cypher to configure for the TLS V1.2 are :
- TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
- TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
- TLS_ECDHE_ECDSA_WITH_AES_256_CCM
- TLS_ECDHE_ECDSA_WITH_AES_128_CCM
- TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
- TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
- TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 
- TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256

The option `HTTP Strict Transport Security` (HSTS) has to be enabled.  

For the configuration of the TLS, the security team will provide 3 certificates : a root certificate, and intermediate certificate and an end entity certificate. The security team will also manage the deployment of the private key into the platform directly.   


