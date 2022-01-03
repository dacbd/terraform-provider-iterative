package utils

import (
	"os"

	"github.com/aohorodnyk/uid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	SynthRegions = map[string]map[string]string{
		"aws": {
			"us-east":  "us-east-1",
			"us-west":  "us-west-1",
			"eu-north": "eu-north-1",
			"eu-west":  "eu-west-1",
		},
		"gcp": {
			"us-east":  "us-east1-c",
			"us-west":  "us-west1-b",
			"eu-north": "europe-north1-a",
			"eu-west":  "europe-west1-d",
		},
		"az": {
			"us-east":  "eastus",
			"us-west":  "westus2",
			"eu-north": "northeurope",
			"eu-west":  "westeurope",
		},
	}
)

func MachinePrefix(d *schema.ResourceData) string {
	prefix := ""
	if _, hasMachine := d.GetOk("machine"); hasMachine {
		prefix = "machine.0."
	}

	return prefix
}

func SetId(d *schema.ResourceData) {
	if len(d.Id()) == 0 {
		d.SetId("iterative-" + uid.NewProvider36Size(8).MustGenerate().String())

		if len(d.Get("name").(string)) == 0 {
			d.Set("name", d.Id())
		}
	}
}

func StripAvailabilityZone(region string) string {
	lastChar := region[len(region)-1]
	if lastChar >= 'a' && lastChar <= 'z' {
		return region[:len(region)-1]
	}
	return region
}

func GetRegion(d *schema.ResourceData) string {
	region := d.Get("region").(string)
	cloud := d.Get("cloud").(string)

	if val, ok := SynthRegions[cloud][region]; ok {
		return val
	}
	if cloud == "aws" {
		return StripAvailabilityZone(region)
	}
	return region
}

func LoadGCPCredentials() string {
	credentialsData := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_DATA")
	if len(credentialsData) == 0 {
		credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if len(credentialsPath) > 0 {
			jsonData, _ := os.ReadFile(credentialsPath)
			credentialsData = string(jsonData)
		}
	}
	return credentialsData
}
