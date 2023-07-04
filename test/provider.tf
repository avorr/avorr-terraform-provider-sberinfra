terraform {
  required_version = ">= 0.14.9"
  required_providers {
    si = {
      source  = "sberbank/devops/si"
#            source  = "cloud/si"
      version = "0.4.7"
    }
  }
  #  backend "s3" {
  #    key                         = "sberdevices/smartapp/asdsadsads.tfstate"
  #    bucket                      = "tfstates"
  #    skip_credentials_validation = true
  #    skip_region_validation      = true
  #    skip_metadata_api_check     = true
  #    force_path_style            = true
  #  }
}

provider si {}
