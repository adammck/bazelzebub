# Bazelzebub

This is just a demo.

```
$ bazel query //example/...
//example:meta
//example:prd-us1-org9876-S1
//example:prd-us1-org9876-S2
//example:prd-us1-org9876-S3  # <--
//example:wf-clu--prd-us1-org9876
//example:wf-dc--prd-eu1
//example:wf-dc--prd-us1
//example:wf-dc--stg-us1
//example:wf-env--prd
//example:wf-env--stg
Loading: 0 packages loaded

$ go build

$ ./bazelzebub example/vars.bzl 
ok

$ bazel query //example/...
//example:meta
//example:prd-us1-org9876-S1
//example:prd-us1-org9876-S2
//example:prd-us1-org9876-S4  # <--
//example:prd-us1-org9876-S5  # <--
//example:wf-clu--prd-us1-org9876
//example:wf-dc--prd-eu1
//example:wf-dc--prd-us1
//example:wf-dc--stg-us1
//example:wf-env--prd
//example:wf-env--stg
```
