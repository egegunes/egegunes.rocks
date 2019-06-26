# egegunes.rocks

It's a little Go web application to get anonymous feedback about myself. There
are no analytics or tracking services included. I collect only your name
(optional) and your message (required).

## Stack

* AWS API Gateway
* AWS Certificate Manager
* AWS Lambda
* AWS DynamoDB
* AWS Route53
* AWS IAM

## Deployment

All you need to run to deploy new version to production is `make`.
