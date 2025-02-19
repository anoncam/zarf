## zarf package deploy

Use to deploy a Zarf package from a local file or URL (runs offline)

### Synopsis

Uses current kubecontext to deploy the packaged tarball onto a k8s cluster.

```
zarf package deploy [PACKAGE] [flags]
```

### Options

```
      --components string    Comma-separated list of components to install.  Adding this flag will skip the init prompts for which components to install
      --confirm              Confirm package deployment without prompting
  -h, --help                 help for deploy
      --insecure --shasum    Skip shasum validation of remote package. Required if deploying a remote package and --shasum is not provided
      --set stringToString   Specify deployment variables to set on the command line (KEY=value) (default [])
      --sget string          Path to public sget key file for remote packages signed via cosign
      --shasum --insecure    Shasum of the package to deploy. Required if deploying a remote package and --insecure is not provided
      --tmpdir string        Specify the temporary directory to use for intermediate files
```

### Options inherited from parent commands

```
  -a, --architecture string   Architecture for OCI images
  -l, --log-level string      Log level when running Zarf. Valid options are: warn, info, debug, trace
      --no-progress           Disable fancy UI progress bars, spinners, logos, etc.
```

### SEE ALSO

* [zarf package](zarf_package.md)	 - Zarf package commands for creating, deploying, and inspecting packages

