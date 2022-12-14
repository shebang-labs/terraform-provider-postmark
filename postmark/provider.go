package postmark

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postmarkSDK "github.com/keighl/postmark"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("POSTMARK_ACCOUNT_TOKEN", nil),
				Description: descriptions["token"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"postmark_server": resourceServer(),
			"postmark_stream": resourceStream(),
			"postmark_domain": resourceDomain(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"postmark_servers": dataSourceServers(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	accountToken := d.Get("account_token").(string)

	var c interface{}
	if accountToken != "" {
		c = postmarkSDK.NewClient("", accountToken)
	}

	if c == nil {
		log.Println("[ERROR] Initializing postmark client is not completed")
		return nil, nil
	}
	log.Println("[INFO] Initializing postmark client")

	return c, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token": "The API account token for postmark API operations.",
	}
}
