package postmark

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postmarkSDK "github.com/keighl/postmark"
)

func dataSourceServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServersRead,
		Schema: map[string]*schema.Schema{
			"totalcount": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"servers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"color": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"apitokens": &schema.Schema{
							Type:      schema.TypeList,
							Computed:  true,
							Sensitive: true,
							MinItems:  1,
							Optional:  true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

type Servers struct {
	TotalCount int                  `json:"TotalCount"`
	Servers    []postmarkSDK.Server `json:"Servers"`
}

func dataSourceServersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := &http.Client{}
	c := m.(*postmarkSDK.Client)

	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", "https://api.postmarkapp.com/servers?count=500&offset=0", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Account-Token", c.AccountToken)
	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	servers := Servers{}
	servers.Servers = make([]postmarkSDK.Server, 0)
	err = json.NewDecoder(r.Body).Decode(&servers)
	newServers := flattenServersData(&servers.Servers)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("servers", newServers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("totalcount", servers.TotalCount); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenServersData(serverItems *[]postmarkSDK.Server) []interface{} {
	if serverItems != nil {
		sis := make([]interface{}, len(*serverItems))

		for i, serverItem := range *serverItems {
			si := make(map[string]interface{})

			si["id"] = serverItem.ID
			si["name"] = serverItem.Name
			si["color"] = serverItem.Color
			si["apitokens"] = serverItem.ApiTokens
			sis[i] = si
		}

		return sis
	}

	return make([]interface{}, 0)
}
