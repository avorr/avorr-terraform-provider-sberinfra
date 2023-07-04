#!/usr/bin/env bash
set -ex
cd ..
echo `date`
name=terraform-provider-si
version=0.4.7
platform=darwin_amd64; goos="darwin"; goarch="amd64"
#platform=darwin_arm64; goos="darwin"; goarch="arm64"
#platform=linux_amd64; goos="linux"; goarch="amd64"
#platform=windows; goos="windows"; goarch="amd64"

provider_dir=${HOME}/.terraform.d/plugins/sberbank/devops/si
binary_dir=${provider_dir}/${version}/${platform}
binary=${name}_v${version}_${platform}

if [[ ${platform} == darwin_arm64 ]]; then
  [ -e ${binary_dir}/${binary} ] && rm ${binary_dir}/${binary}
fi

mkdir -p ${binary_dir}
go version
go mod tidy -v
#go mod vendor -v
#go build -mod=vendor -v -o ${binary}
GOOS=${goos} GOARCH=${goarch} go build -v -o ${binary}
mv ${binary} ${binary_dir}/${binary}
cd test

rm -rf .terraform/ || true
rm ./.terraform.lock.hcl || true

export SI_TIMEOUT=600
#export TF_LOG=DEBUG
#export TF_LOG=INFO
#export TF_LOG=ERROR

terraform init

#terraform plan
#terraform apply -auto-approve
#terraform destroy -auto-approve

#terraform import si_vdc.vdc c41a6b76-ddfe-4d49-a762-ea659becf35f
#terraform state show -no-color si_vdc.vdc >> project.tf
#terraform import si_vm.vm e94bea8e-7ea3-49da-b91f-0c71092da6ff