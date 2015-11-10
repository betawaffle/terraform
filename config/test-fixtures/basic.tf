variable "foo" {
    default = "bar"
    description = "bar"
}

provider "aws" {
  access_key = "foo"
  secret_key = "bar"
}

provider "do" {
  api_key = "${var.foo}"
}

resource "aws_security_group" "firewall" {
    count = 5
}

resource aws_instance "web" {
    ami = "${var.foo}"
    security_groups = [
        "foo",
        "${aws_security_group.firewall.foo}"
    ]

    network_interface {
        device_index = 0
        description = "Main network interface"
    }

    provisioner "file" {
        source = "foo"
        destination = "bar"
    }
}

resource "aws_instance" "db" {
    security_groups = "${aws_security_group.firewall.*.id}"
    VPC = "foo"

    depends_on = ["aws_instance.web"]

    provisioner "file" {
        source = "foo"
        destination = "bar"
    }
}

resource "aws_iam_policy" "policy" {
    name = "test_policy"
    path = "/"
    description = "My test policy"
    policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

output "web_ip" {
    value = "${aws_instance.web.private_ip}"
}

atlas {
    name = "mitchellh/foo"
}
