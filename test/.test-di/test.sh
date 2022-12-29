#!/usr/bin/env bash

rm -rf .terraform/
rm terraform.tfstate*
rm ./.terraform.lock.hcl

export TF_LOG=DEBUG
#export TF_LOG=INFO
#export TF_LOG=ERROR

terraform init
#./imports.sh
terraform plan
#terraform apply
#terraform apply -auto-approve


# stand: nlpf:
#1 ignite (in mem) cluster 2
#1 openshift
# instances: alpha, sigma/ift,pt,uat,prod

# stand:
#