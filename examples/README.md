# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI or OpenTofu.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* `provider/provider.tf` example file for the provider index page
* `data-sources/full data source name/data-source.tf` example file for the named data source page
* `resources/full resource name/resource.tf` example file for the named data source page

Before running the examples, you need to set a project and IAM token:

```shell
export MWS_PROJECT=<your_project>
export MWS_TOKEN=$(mws iam create-token)
```

You can also set the following parameters in the provider configuration:

```hcl
provider "mws" {
  project   = "<your_project>"
  mws_token = "<your_iam_token>"
}
```
