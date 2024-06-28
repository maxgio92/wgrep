## wgrep

Wgrep print lines that match patterns in web pages

```
wgrep PATTERN URL [flags]
```

### Options

```
      --async                               Whether to scrape with asynchronous jobs. (default true)
      --connection-pool-size int            The maximum number of idle connections across all hosts. (default 1000)
      --connection-pool-size-per-host int   The maximum number of idle connections across for each host. (default 1000)
      --connection-timeout int              The maximum amount of time in milliseconds a dial will wait for a connect to complete. (default 180000)
      --element string                      The HTLM element to filer on while searching for the pattern (default "p")
  -h, --help                                help for wgrep
      --idle-connection-timeout int         The maximum amount of time in milliseconds a connection will remain idle before closing itself. (default 120000)
  -i, --ignore-case                         Whether to search for the pattern case insensitive
      --include string                      Search only pages whose URL matches teh include regular expression.
      --keep-alive-interval int             The interval between keep-alive probes for an active network connection. (default 30000)
      --max-body-size int                   The maximum size in bytes a response body is read for each request. (default 524288)
  -r, --recursive                           Inspect all web pages recursively by following each hypertext reference.
      --tls-handshake-timeout int           The maximum amount of time in milliseconds a connection will wait for a TLS handshake. (default 30000)
  -v, --verbose                             Enable verbosity to log all visited HTTP(s) files
```

