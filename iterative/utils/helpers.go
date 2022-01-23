package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aohorodnyk/uid"
	"github.com/blang/semver/v4"
	"github.com/google/go-github/v42/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetCML(version string) string {
	// "v"semver
	if strings.HasPrefix(version, "v") {
		ver, err := semver.Make(version[1:])
		if err != nil {
			return getSemverCML(ver)
		}
	}
	// semver
	ver, err := semver.Make(version)
	if err != nil {
		return getSemverCML(ver)
	}
	// npm install string
	if version != "" {
		return getNPMCML(version)
	}
	//default latest
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "iterative", "cml")
	if err != nil {
		for _, asset := range release.Assets {
			if *asset.Name == "cml-linux" {
				return getGHCML(*asset.BrowserDownloadURL)
			}
		}
	}
	// original fallback
	return getNPMCML("@dvcorg/cml")
}
func getGHCML(v string) string {
	ghCML := "curl %s -o /bin/cml && chmod +x /bin/cml"
	return fmt.Sprint(ghCML, v)
}
func getNPMCML(v string) string {
	npmCML := "sudo npm config set user 0 && sudo npm install --global %s"
	return fmt.Sprint(npmCML, v)
}
func getSemverCML(sv semver.Version) string {
	directDownloadVersion, _ := semver.ParseRange(">=0.10.0")
	if directDownloadVersion(sv) {
		client := github.NewClient(nil)
		release, _, err := client.Repositories.GetReleaseByTag(context.Background(), "iterative", "cml", "v"+sv.String())
		if err != nil {
			for _, asset := range release.Assets {
				if *asset.Name == "cml-linux" {
					return getGHCML(*asset.BrowserDownloadURL)
				}
			}
		}
	}
	// npm install
	return getNPMCML("@dvcorg/cml@v" + sv.String())

}

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
