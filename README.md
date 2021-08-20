# citizenlab_translator
Fetch and translate `github.com/CitizenLabDotCo/citizenlab` source language files using the [AWS Amazon Translate](https://aws.amazon.com/translate/) service.

### How to
1. make sure your AWS credentials are in place (see the [official documentation](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials) for more info)

1. `go run . --lang it --output ./output --type yml`

### Limitations
* as per [github.com/CitizenLabDotCo/citizenlab/issues/300](https://github.com/CitizenLabDotCo/citizenlab/issues/300) we are currently only supporting 2 input files

## TODO
1. handle case where an item is removed from source
1. remove `<nil>` from yml output