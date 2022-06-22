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
"<token_creation>","2022-05-11T14:37:17.037172945Z","foo/prod/","auth/approle/login","approle","batch"
"<token_creation>","2022-05-11T14:37:18.493810518Z","foo/dev/","auth/approle/login","approle","batch"
"<token_creation>","2022-05-11T14:37:21.202577225Z","foo/prod/","auth/approle/login","approle","batch"
"<token_creation>","2022-05-11T14:37:24.388627203Z","foo/dev/","auth/approle/login","approle","batch"
"<token_creation>","2022-05-11T14:37:29.684249077Z","foo/qa/","auth/approle/login","approle","batch"
...

$ tail -f audit.log | ./revizio -errors
"<error>","2022-05-11T14:37:17.037172945Z","foo/","sys/internal/ui/mounts/secret/wif-adfs","ns_system","","permission denied"
"<error>","2022-05-11T14:37:18.493810518Z","<root>","sys/internal/ui/resultant-acl","system","","1 error occurred:\n\t* permission denied\n\n"
"<error>","2022-05-11T14:37:21.202577225Z","foo/","auth/kubernetes/login","kubernetes","","service account name not authorized"
"<error>","2022-05-11T14:37:24.388627203Z","<root>","sys/internal/ui/resultant-acl","system","","1 error occurred:\n\t* permission denied\n\n"
"<error>","2022-05-11T14:37:29.684249077Z","bar/","auth/approle/role/gcp_pipeline/secret-id","approle","","1 error occurred:\n\t* permission denied\n\n"
...

$ tail -f audit.log | ./revizio -errors -tokens
"<token_creation>","2022-05-11T14:37:17.037172945Z","foo/prod/","auth/approle/login","approle","batch"
"<token_creation>","2022-05-11T14:37:18.493810518Z","foo/dev/","auth/approle/login","approle","batch"
"<token_creation>","2022-05-11T14:37:20.748267900Z","foo/prod/","auth/approle/login","approle","batch"
"<error>","2022-05-11T14:37:21.202577225Z","bar/","sys/internal/ui/mounts/secret/abcd","ns_system","","permission denied"
"<error>","2022-05-11T14:37:24.388627203Z","bar/","sys/internal/ui/mounts/secret/abcd","ns_system","","permission denied"
"<token_creation>","2022-05-11T14:37:29.684249077Z","foo/qa/","auth/approle/login","approle","batch"
"<token_creation>","2022-05-11T14:37:35.162773048Z","foo/prod/","auth/approle/login","approle","batch"
...
```

## TODO
- Flesh out summary view
- Bring in time, latency
