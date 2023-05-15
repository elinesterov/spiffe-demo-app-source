# SPIFFE Demo App


## Setup

This project uses `ko` to build images. [You can learn more about ko here](https://ko.build/).

Install KO
```
brew install ko
```

ko will push images to the repository defined by KO\_DOCKER\_REPO.

You can use the `publish_poc` target to push to the PoC ECR repository

As example to pushlish to ECR PoC
```
aws ecr get-login-password --region us-west-2 | \
    ko login --username AWS --password-stdin 771189981606.dkr.ecr.us-west-2.amazonaws.com
```

