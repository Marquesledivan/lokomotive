---
title: lokoctl health
weight: 10
---

Get the health of a cluster

### Synopsis

Get the health of a cluster

```
lokoctl health [flags]
```

### Options

```
  -h, --help   help for health
```

### Options inherited from parent commands

```
      --kubeconfig-file string   Path to a kubeconfig file. If empty, the following precedence order is used:
                                   1. Cluster asset dir when a lokocfg file is present in the current directory.
                                   2. KUBECONFIG environment variable.
                                   3. ~/.kube/config file.
      --lokocfg string           Path to lokocfg directory or file (default "./")
      --lokocfg-vars string      Path to lokocfg.vars file (default "./lokocfg.vars")
```

### SEE ALSO

* [lokoctl](lokoctl.md)	 - Manage Lokomotive clusters

