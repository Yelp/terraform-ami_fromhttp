package yelpaws

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/hashicorp/terraform/builtin/providers/aws"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const jenkinsUrlFmt = "https://jenkins.yelpcorp.com/job/promote-generic%s%s-ami/lastSuccessfulBuild/artifact/"
const artifactFmt = "account-%s-aws_region-%s_ami_id.txt"

var pvOnly = map[string]bool{"m1": true, "m2": true, "c1": true, "t1": true}

func Provider() terraform.ResourceProvider {
	provider := aws.Provider().(*schema.Provider)
	instance := provider.ResourcesMap["aws_instance"]
	ami := instance.Schema["ami"]
	ami.Required = false
	ami.Optional = true
	ami.Computed = true
	instance.Create = wrapCreate(instance.Create)
	return provider
}

func wrapCreate(create func(*schema.ResourceData, interface{}) error) func(*schema.ResourceData, interface{}) error {
	wrapped := func(d *schema.ResourceData, meta interface{}) error {
		if _, ok := d.GetOk("ami"); !ok {
			ami := getAMI("lucid", d.Get("instance_type").(string))
			d.Set("ami", ami)
		}
		return create(d, meta)
	}
	return wrapped
}

func buildAmiUrl(osdist string, hvm bool, account string, region string) string {
	var hvmSuffix string
	if hvm {
		hvmSuffix = "-hvm"
	}
	return fmt.Sprintf(jenkinsUrlFmt, osdist, hvmSuffix) + fmt.Sprintf(artifactFmt, account, region)
}

func getAMI(osdist string, instanceType string) string {
	// TODO: get account, region from provider config
	url := buildAmiUrl(osdist, supportsHVM(instanceType), "dev", "us-west-1")
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	ami_id, _ := ioutil.ReadAll(resp.Body)
	return strings.TrimSpace(string(ami_id))
}

func supportsHVM(instanceType string) bool {
	instanceClass := instanceType[:2]
	return !pvOnly[instanceClass]
}
