packer {
  required_plugins {
    amazon = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "region" {}
variable "ami_name" {}
variable "ami_users" {}
variable "instance_type" {}
variable "profile" {}
variable "subnet_id" {}
variable "source_ami_name" {}
variable "source_ami_device_type" {}
variable "source_ami_virtualization_type" {}
variable "source_ami_owners" {}
variable "ssh_username" {}
variable "device_name" {}
variable "volume_size" {}
variable volume_type {}
variable "file_paths" {}
variable destination_path {}
variable script_path {}
variable "shell_environment_vars" {}

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

source "amazon-ebs" "debian_ami" {
  ami_name      = "${var.ami_name}-${local.timestamp}"
  ami_users     = var.ami_users
  instance_type = var.instance_type
  profile       = var.profile
  region        = var.region
  subnet_id     = var.subnet_id
  source_ami_filter {
    filters = {
      name                = var.source_ami_name
      root-device-type    = var.source_ami_device_type
      virtualization-type = var.source_ami_virtualization_type
    }
    most_recent = true
    owners      = var.source_ami_owners
  }
  ssh_username = var.ssh_username

  launch_block_device_mappings {
    device_name           = var.device_name
    delete_on_termination = true
    volume_size           = var.volume_size
    volume_type           = var.volume_type
  }
}

build {
  name    = "build-packer"
  sources = ["source.amazon-ebs.debian_ami"]

  provisioner "file" {
    sources     = var.file_paths
    destination = var.destination_path
  }

  provisioner "shell" {
    environment_vars = var.shell_environment_vars
    script           = var.script_path
  }

  post-processor "manifest" {
    output     = "manifest.json"
    strip_path = true
  }
}
