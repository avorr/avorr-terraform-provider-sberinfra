resource "si_security_group" "iam" {
  name   = "iam"
  vdc_id = si_vdc.vdc.id
  security_rule {
    ethertype   = "IPv4"
    direction   = "ingress"
    protocol    = "tcp"
    cidr_prefix = "172.21.21.10/28"
    from_port   = 443
    to_port     = 444
  }
  security_rule {
    ethertype   = "IPv4"
    direction   = "ingress"
    protocol    = "tcp"
    cidr_prefix = "172.21.21.10/28"
    from_port   = 80
    to_port     = 80
  }
}

resource "si_security_group" "kafka" {
  name   = "kafka"
  vdc_id = si_vdc.vdc.id
  security_rule {
    ethertype = "IPv4"
    direction = "ingress"
    protocol  = "tcp"
    from_port = 9092
    to_port   = 9092
  }
  security_rule {
    ethertype = "IPv4"
    direction = "ingress"
    protocol  = "tcp"
    from_port = 2181
    to_port   = 2181
  }
}

#resource "si_security_group" "import" {}