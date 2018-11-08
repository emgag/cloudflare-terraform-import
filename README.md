# quick'n dirty cloudflare to terraform importer for DNS records

(Only tested with A, CAA, CNAME, MX, TXT records so far) 

## build & run

Requires Go >= 1.11

```
$ git clone git@github.com:emgag/cloudflare-terraform-import.git
$ go get && make
$ export CF_API_EMAIL="user@example.org" 
$ export CF_API_KEY="..."
$ ./cti import example.org
``` 

* Copy contents of `import.tf` to your terraform config.
* Review 
* Run `import.sh` to import records into state
