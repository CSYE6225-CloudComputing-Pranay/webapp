### AWS configuration variables
region    = "us-east-1"
subnet_id = "subnet-0ba56c90b74d73fd5"

### Variables for EC2 instance creation
file_paths             = ["./assessment-application", "./build/package/users.csv", "./build/package/assessment-application.service", "./build/package/amazon-cloudwatch-agent.json"]
script_path            = "./build/package/setup.sh"
destination_path       = "/tmp/"
ssh_username           = "admin"
shell_environment_vars = ["DEBIAN_FRONTEND=noninteractive", "CHECKPOINT_DISABLE=1"]

### Variables for AMI creation
ami_name      = "csye6225-debian-instance-ami"
ami_users     = ["608008135379"]
instance_type = "t2.micro"
profile       = "default"
volume_size   = 25
device_name   = "/dev/xvda"
volume_type   = "gp2"

### Filters to fetch source AMI
source_ami_name                = "debian-12-amd64-*"
source_ami_device_type         = "ebs"
source_ami_virtualization_type = "hvm"
source_ami_owners              = ["amazon"]