#!/usr/bin/env bash
#source ~/.bash_profile
set -ex
cd ..
echo `date`
name=terraform-provider-di
version=0.3.13
platform=darwin_amd64
#platform=linux_amd64

provider_dir=${HOME}/.terraform.d/plugins/sberbank/devops/di
binary_dir=${provider_dir}/${version}/${platform}
binary=${name}_v${version}_${platform}

mkdir -p ${binary_dir}
go version
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

#cd test/
cd tf-files-example/
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
#terraform apply
ls -l ~/.terraform.d/plugins/sberbank/devops/di/0.3.13/darwin_amd64/terraform-provider-di_v0.3.13_darwin_amd64
#terraform apply -auto-approve
#terraform destroy -auto-approve

#terraform import di_patroni.patroni 878b50c7-5eb4-4d05-a81a-76597c2ddb84
#terraform import di_project.terraformtest2 c41a6b76-ddfe-4d49-a762-ea659becf35f
#terraform import di_project.nlpf_test 41197683-fb3e-42a2-acb7-25089367e9d5
#terraform state show -no-color di_project.terraformtest2
#terraform state show -no-color di_project.terraformtest2 >> project.tf

#terraform import di_vm.vm1 e94bea8e-7ea3-49da-b91f-0c71092da6ff
#terraform import di_vm.server_17756439_24_06_2022_165429 79d2ef40-46ea-4415-852a-f5a456728f5a

#${binary_dir}/${binary} import 9efde59e-2db8-4b46-ade8-8cd952890e53
#${binary_dir}/${binary} import 82f0c513-08ce-4b74-9b30-5eb287c895a9
#${binary_dir}/${binary} import f661f5fe-c436-470b-912a-7a21e350920a
#bash imports.sh
