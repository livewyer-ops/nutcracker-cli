# Nutcracker-cli

Commandline interface for github.com/nutmegdevelopment/nutcracker

This is designed to be used primarily by automated services, and such has very basic functionality.
More advanced functionality will be in the upcoming github.com/nutmegdevelopment/nutcracker-ui project.

Nutcracker-cli requires three arguments:
```
   --server, -s 	Nutcracker server.  e.g localhost:443
   --id, -i 		Nutcracker API ID
   --key, -k 		Nutcracker API key
```

Supported commands are list and get.  Get requires an additional parameter, --name or -n which is the name of the secret to get.

Full syntax:

```
nutcracker-cli -s 1.2.3.4:8443 -i abcd -k xyz123 list
nutcracker-cli -s 1.2.3.4:8443 -i abcd -k xyz123 get -n foo
```
