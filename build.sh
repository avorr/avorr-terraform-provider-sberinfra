#!/usr/bin/env bash
set -ex

echo `date`
name=terraform-provider-di
version=0.3.12
platform=darwin_amd64

provider_dir=${HOME}/.terraform.d/plugins/sberbank/devops/di
binary_dir=${provider_dir}/${version}/${platform}
binary=${name}_v${version}_${platform}

mkdir -p ${binary_dir}
go mod tidy -v
#go mod vendor -v
#go build -mod=vendor -v -o ${binary}
go build -v -o ${binary}
#./${binary} import
#rm test/.terraform.lock.hcl || true
cp ${binary} ${binary_dir}/${binary}
#mkdir -p test/.terraform/plugins/sberbank/devops/di/${version}/${platform}/
#cp ${binary} test/.terraform/plugins/sberbank/devops/di/${version}/${platform}/
rm ${binary}

cd test/
rm -rf .terraform/ || true
#rm terraform.tfstate* || true
rm ./.terraform.lock.hcl || true
#rm ./inventory.bin || true

export TF_LOG=DEBUG
export DI_TIMEOUT=7000
#export TF_LOG=INFO
#export TF_LOG=ERROR
export DI_ANSIBLE_PASSWORD=False
export INVENTORY_DISABLE=False

terraform init
#./imports.sh
#terraform plan
terraform apply
#terraform apply -auto-approve

#terraform import di_patroni.patroni 878b50c7-5eb4-4d05-a81a-76597c2ddb84
