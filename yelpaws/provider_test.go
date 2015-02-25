package yelpaws

import (
	"testing"
)

func TestSupportsHVM(t *testing.T) {
	cases := []struct{
		in string
		expected bool
	} {
		{"t1.micro", false},
		{"c1.large", false},
		{"c4.2xlarge", true},
		{"t2.micro", true},
	}

	for _, c := range cases {
		got := supportsHVM(c.in)
		if got != c.expected {
			t.Errorf("supportsHVM(%q) == %q, expected %t", c.in, got, c.expected)
		}
	}
}

func TestBuildAmiUrl(t *testing.T) {
	cases := []struct{
		os string
		hvm bool
		account string
		region string
		expected string
	} {
		{"lucid", false, "dev", "us-west-1", "https://jenkins.yelpcorp.com/job/promote-genericlucid-ami/lastSuccessfulBuild/artifact/account-dev-aws_region-us-west-1_ami_id.txt"},
		{"lucid", true, "dev", "us-west-1", "https://jenkins.yelpcorp.com/job/promote-genericlucid-hvm-ami/lastSuccessfulBuild/artifact/account-dev-aws_region-us-west-1_ami_id.txt"},
		{"testbuntu", false, "awstest", "us-test-1", "https://jenkins.yelpcorp.com/job/promote-generictestbuntu-ami/lastSuccessfulBuild/artifact/account-awstest-aws_region-us-test-1_ami_id.txt"},
	}

	for _, c := range cases {
		got := buildAmiUrl(c.os, c.hvm, c.account, c.region)
		if got != c.expected {
			t.Errorf("buildAmiUrl(%q, %t, %q, %q) == %q, expected %q", c.os, c.hvm, c.account, c.region, got, c.expected)
		}
	}
}
