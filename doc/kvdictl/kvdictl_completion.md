## kvdictl completion

Generate completion script

### Synopsis

To load completions:

Bash:

  $ source <(kvdictl completion bash)

To load completions for each session, execute once:

Linux:
  $ kvdictl completion bash > /etc/bash_completion.d/kvdictl
MacOS:
  $ kvdictl completion bash > /usr/local/etc/bash_completion.d/kvdictl

Zsh:

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ kvdictl completion zsh > "${fpath[1]}/_kvdictl"

  You will need to start a new shell for this setup to take effect.

Fish:

  $ kvdictl completion fish | source
  # To load completions for each session, execute once:
  $ kvdictl completion fish > ~/.config/fish/completions/kvdictl.fish


```
kvdictl completion [bash|zsh|fish|powershell]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
  -C, --ca-file string         the CA certificate to use to verify the API certificate
  -c, --config string          configuration file (default "$HOME/.kvdi.yaml")
  -f, --filter string          a jmespath expression for filtering results (where applicable)
  -k, --insecure-skip-verify   skip verification of the API server certificate
  -o, --output string          the format to dump results in (default "json")
  -s, --server string          the address to the kvdi API server (default "https://127.0.0.1")
  -u, --user string            the username to use when authenticating against the API (default "admin")
```

### SEE ALSO

* [kvdictl](kvdictl.md)	 - 

###### Auto generated by spf13/cobra on 5-Mar-2021