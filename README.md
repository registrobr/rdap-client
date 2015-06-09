# rdap-client

RDAP client usage

```
NAME:
   rdap - RDAP cgclient

USAGE:
   rdap [global options] OBJECT

VERSION:
   0.0.1

AUTHOR(S): 
   NIC.br 

GLOBAL OPTIONS:
   --cache "/home/joao/.rdap"                          directory for caching bootstrap and RDAP data
   --bootstrap "https://data.iana.org/rdap/%s.json"    RDAP bootstrap service URL
   --no-cache                                          don't cache anything
   --skip-tls-verification, -S                         skip TLS verification
   --domain                                            force query for a domain object
   --asn                                               force query for an ASN object
   --ip                                                force query for an IP object
   --ipnetwork                                         force query for an IP Network object
   --entity                                            force query for an Entity object
   --host, -H                                          host where to send the query (bypass bootstrap)
   --help, -h                                          show help
   --version, -v                                       print the version

```
