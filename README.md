# revizio

```sh
$ tail -f audit.log | ./revizio # Outputs CSV of token creation information
...

$ tail -f audit.log | ./revizio -detailed # Provides additional info including request_id, requesting token, token_ttl, etc.
...

$ cat audit.log | ./revizio -summary # Provides some category totals from a preexisting audit log
```