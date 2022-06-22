# revizio

## Usage
```sh
$ tail -f audit.log | ./revizio

  -buffersize int
    	Adjust line buffer size max. (default 262144)
  -errors
    	Output errors.
  -summary
    	Summarized output instead of CSV.
  -tokens
    	Output token creation metadata.
  -verbose
    	More detail.
...

$ tail -f audit.log | ./revizio -tokens
"<token_creation>","foo/prod/","auth/approle/login","approle","batch"
"<token_creation>","foo/dev/","auth/approle/login","approle","batch"
"<token_creation>","foo/prod/","auth/approle/login","approle","batch"
"<token_creation>","foo/dev/","auth/approle/login","approle","batch"
"<token_creation>","foo/qa/","auth/approle/login","approle","batch"
...

$ tail -f audit.log | ./revizio -errors
"<error>","foo/","sys/internal/ui/mounts/secret/wif-adfs","ns_system","","permission denied"
"<error>","<root>","sys/internal/ui/resultant-acl","system","","1 error occurred:\n\t* permission denied\n\n"
"<error>","foo/","auth/kubernetes/login","kubernetes","","service account name not authorized"
"<error>","<root>","sys/internal/ui/resultant-acl","system","","1 error occurred:\n\t* permission denied\n\n"
"<error>","bar/","auth/approle/role/gcp_pipeline/secret-id","approle","","1 error occurred:\n\t* permission denied\n\n"
...

$ tail -f audit.log | ./revizio -errors -tokens
"<token_creation>","foo/prod/","auth/approle/login","approle","batch"
"<token_creation>","foo/dev/","auth/approle/login","approle","batch"
"<token_creation>","foo/prod/","auth/approle/login","approle","batch"
"<error>","bar/","sys/internal/ui/mounts/secret/abcd","ns_system","","permission denied"
"<error>","bar/","sys/internal/ui/mounts/secret/abcd","ns_system","","permission denied"
"<token_creation>","foo/qa/","auth/approle/login","approle","batch"
"<token_creation>","foo/prod/","auth/approle/login","approle","batch"
...

```

## TODO
- Flesh out summary view
- Bring in time, latency
