# Securitytxt service

> Reference:
> * https://securitytxt.org/
> * https://github.com/securitytxt/security-txt
> * https://www.rfc-editor.org/info/rfc9116

This is a small go service that produces a securitytxt file from environment variables and serving it through a http server

The purpose for the container is to deploy in a k8s, or similar, context and adding it to the ingress on one of the following paths
`/.well-known/security.txt` or `/security.txt`

## Example

```bash 
docker run -i \
    --env "CONTACT_URIS=mailto:security@example.com tel:+012-3456-789" \
    --env "EXPIRES=$(date --date="+1 year" +%Y-%m-%d)" \
    --publish "8080:8080" \
    modfin/securitytxt:latest 
```



## Environment variables
```
COMMENT                (string)
EXPIRES                (yyyy-mm-dd or RFC3999)
CONTACT_URIS           ([]string, del " ")
ACKNOWLEDGMENT_URIS    ([]string, del " ")
CANONICAL_URIS         ([]string, del " ")
ENCRYPTION_URIS        ([]string, del " ")
HIRING_URIS            ([]string, del " ")
PREFERRED_LANGUAGES    ([]string, del " ")
POLICY_URIS            ([]string, del " ")

RAW_SECURITY_TXT       (string, overrides all other values)
```
