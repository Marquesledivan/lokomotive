---
title: lokoctl cluster destroy
weight: 10
---

Destroy a cluster

### Synopsis

Destroy a cluster

```
lokoctl cluster destroy [flags]
```

### Options

```
      --confirm   Destroy cluster without asking for confirmation
  -h, --help      help for destroy
  -v, --verbose   Show output from Terraform
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

* [lokoctl cluster](lokoctl_cluster.md)	 - Manage a cluster

