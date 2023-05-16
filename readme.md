# SPIFFE Demo App


## Setup

This project uses `ko` to build images. [You can learn more about ko here](https://ko.build/).

Install KO
```
brew install ko
```

ko will push images to the repository defined by KO\_DOCKER\_REPO.

You can use the `publish_poc` target to push to the PoC ECR repository

```
make publish_poc
```

Or if you want to generate the manifest to deploy the build

```
make resolve_poc
```
