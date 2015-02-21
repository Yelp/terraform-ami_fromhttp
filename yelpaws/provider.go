package yelpaws

import (
	"github.com/hashicorp/terraform/builtin/providers/aws"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

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

func getAMI(lsbdist string, instanceType string) string {
	// TODO: Return an apporpriate AMI from jenkins
	return "ami-39501209"
}

func supportsHVM(instanceType string) bool {
	instanceClass := instanceType[0:1]
	return !pvOnly[instanceClass]
}
