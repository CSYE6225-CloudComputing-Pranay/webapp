packer {
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "ami_name" {
  type    = string
  default = "csye6225-debian-instance-ami"
}

variable "demo_account_id" {
  type    = string
  default = "608008135379"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

variable "profile" {
  type    = string
  default = "default"
}

variable "region" {
  type    = string
  default = "us-east-1"
}

variable "subnet_id" {
  type    = string
  default = "subnet-0ba56c90b74d73fd5"
}

variable "source_ami_name" {
  type    = string
  default = "debian-12-amd64-*"
}

variable "source_ami_device_type" {
  type    = string
  default = "ebs"
}

variable "source_ami_virtualization_type" {
  type    = string
  default = "hvm"
}

variable "source_ami_owners" {
  type    = list(string)
  default = ["amazon"]
}

locals {
  timestamp       = regex_replace(timestamp(), "[- TZ:]", "")
  demo_account_id = "081235755261"
}

source "amazon-ebs" "debian_ami" {
  ami_name      = "${var.ami_name}-${local.timestamp}"
  ami_users     = ["${var.demo_account_id}"]
  instance_type = "${var.instance_type}"
  profile       = "${var.profile}"
  region        = "${var.region}"
  subnet_id     = "${var.subnet_id}"
  source_ami_filter {
    filters = {
      name                = "${var.source_ami_name}"
      root-device-type    = "${var.source_ami_device_type}"
      virtualization-type = "${var.source_ami_virtualization_type}"
    }
    most_recent = true
    owners      = "${var.source_ami_owners}"
  }
  ssh_username = "admin"

  launch_block_device_mappings {
    device_name           = "/dev/xvda"
    delete_on_termination = true
    volume_size           = 25
    volume_type           = "gp2"
  }
}

build {
  name = "build-packer"
  sources = [
    "source.amazon-ebs.debian_ami"
  ]

  provisioner "file" {
    sources = [
      "./assessment-application",
      "./build/package/users.csv"
    ]
    destination = "/tmp/"
  }

  provisioner "shell" {
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1"
    ]
    script = "./build/package/setup.sh"
  }
}

