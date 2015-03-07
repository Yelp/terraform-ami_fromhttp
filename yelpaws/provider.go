package yelpaws

import (
	"fmt"
	"github.com/hashicorp/terraform/builtin/providers/aws"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"net/http"
	"strings"
)

const jenkinsUrlFmt = "https://jenkins.yelpcorp.com/job/promote-generic%s%s-ami/lastSuccessfulBuild/artifact/"
const artifactFmt = "account-%s-aws_region-%s_ami_id.txt"

var pvOnly = map[string]bool{"m1": true, "m2": true, "c1": true, "t1": true}

func Provider() terraform.ResourceProvider {
	provider := aws.Provider().(*schema.Provider)
	awsConfig := provider.Schema
	awsConfig["account"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		Description: "The Yelp AWS account that is being used.\n" +
			"Current known account are (dev, devb, stagea, stageb, prod)",
		InputDefault: "dev",
	}
	instanceResource := provider.ResourcesMap["aws_instance"]
	ami := instanceResource.Schema["ami"]
	ami.Required = false
	ami.Optional = true
	ami.Computed = true
	instanceResource.Schema["ubuntu"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "lucid",
		ForceNew: true,
	}
	instanceResource.Create = wrapCreate(instanceResource.Create)
	instanceResource.Read = wrapCRUD(instanceResource.Read)
	instanceResource.Update = wrapCRUD(instanceResource.Update)
	instanceResource.Delete = wrapCRUD(instanceResource.Delete)

	return &schema.Provider{
		Schema: awsConfig,
		ResourcesMap: map[string]*schema.Resource{
			"yelpaws_instance": instanceResource,
		},
		ConfigureFunc: wrapConfigure(provider.ConfigureFunc),
	}
}

// TODO: Replace region/account with an AMIFetcher interface
type YelpAWSClient struct {
	Account string
	Region  string
	*aws.AWSClient
}

func wrapConfigure(configure func(*schema.ResourceData) (interface{}, error)) func(*schema.ResourceData) (interface{}, error) {
	wrapped := func(d *schema.ResourceData) (interface{}, error) {
		awsClient, error := configure(d)
		if error != nil {
			return awsClient, error
		}
		client := YelpAWSClient{
			Account:   d.Get("account").(string),
			Region:    d.Get("region").(string),
			AWSClient: awsClient.(*aws.AWSClient),
		}
		return &client, nil
	}
	return wrapped
}

func wrapCRUD(f func(*schema.ResourceData, interface{}) error) func(*schema.ResourceData, interface{}) error {
	wrapped := func(d *schema.ResourceData, meta interface{}) error {
		yelpClient := meta.(*YelpAWSClient)
		return f(d, yelpClient.AWSClient)
	}
	return wrapped
}

func wrapCreate(create func(*schema.ResourceData, interface{}) error) func(*schema.ResourceData, interface{}) error {
	wrapped := func(d *schema.ResourceData, meta interface{}) error {
		yelpClient := meta.(*YelpAWSClient)
		if _, ok := d.GetOk("ami"); !ok {
			ami := getAMI(
				d.Get("ubuntu").(string),
				d.Get("instance_type").(string),
				yelpClient.Account,
				yelpClient.Region,
			)
			d.Set("ami", ami)
		}
		return create(d, yelpClient.AWSClient)
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

func getAMI(osdist string, instanceType string, account string, region string) string {
	url := buildAmiUrl(osdist, supportsHVM(instanceType), account, region)
	resp, _ := http.Get(url)
	// TODO: Smart things on error
	if resp.StatusCode != 200 {
		return "bad-ami"
	}
	defer resp.Body.Close()
	ami_id, _ := ioutil.ReadAll(resp.Body)
	return strings.TrimSpace(string(ami_id))
}

func supportsHVM(instanceType string) bool {
	instanceClass := instanceType[:2]
	return !pvOnly[instanceClass]
}
