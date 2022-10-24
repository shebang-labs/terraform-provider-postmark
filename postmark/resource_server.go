package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postmarkSDK "github.com/keighl/postmark"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,
		Schema: map[string]*schema.Schema{
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
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Acknowledgment for the implementation for this struct: https://github.com/keighl/postmark/blob/master/servers.go
// Server represents a server registered in your Postmark account
type Server struct {
	// ID of server
	ID int64
	// Name of server
	Name string
	// ApiTokens associated with server.
	ApiTokens []string
	// Color of the server in the rack screen. Purple Blue Turquoise Green Red Yellow Grey
	Color string
	// Delivery type of server
	DeliveryType string
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{}
	c := m.(*postmarkSDK.Client)

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/servers", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Account-Token", c.AccountToken)

	server := Server{}

	server.Name = d.Get("name").(string)
	server.Color = d.Get("color").(string)
	server.DeliveryType = d.Get("delivery_type").(string)
	if server.DeliveryType != "live" && server.DeliveryType != "Sandbox" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Postmark server",
			Detail:   "delivery_type must be either live or Sandbox",
		})
		return diags
	}
	body, err := json.Marshal(server)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	if err != nil {
		return diag.FromErr(err)
	}
	err = json.NewDecoder(res.Body).Decode(&server)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(strconv.FormatInt(server.ID, 10))
	d.Set("apitokens", flattenStringList(server.ApiTokens))

	return diags
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := &http.Client{}
	c := m.(*postmarkSDK.Client)

	if d.HasChange("name") || d.HasChange("color") {
		serverId := d.Id()
		server := postmarkSDK.Server{}
		server.Name = d.Get("name").(string)
		server.Color = d.Get("color").(string)
		req, err := http.NewRequest("PUT", "https://api.postmarkapp.com/servers/"+serverId, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Postmark-Account-Token", c.AccountToken)
		if err != nil {
			return diag.FromErr(err)
		}
		body, err := json.Marshal(server)
		if err != nil {
			return diag.FromErr(err)
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		res, err := client.Do(req)
		if err != nil {
			return diag.FromErr(err)
		}
		defer res.Body.Close()

		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceServerRead(ctx, d, m)
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*postmarkSDK.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	server, err := c.GetServer(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", server.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("color", server.Color); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("apitokens", flattenStringList(server.ApiTokens)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{}
	c := m.(*postmarkSDK.Client)

	serverId := d.Id()
	req, err := http.NewRequest("DELETE", "https://api.postmarkapp.com/servers/"+serverId, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Account-Token", c.AccountToken)
	if err != nil {
		return diag.FromErr(err)
	}
	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func flattenStringList(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}
