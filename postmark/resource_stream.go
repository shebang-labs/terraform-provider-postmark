package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStream() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStreamCreate,
		ReadContext:   resourceStreamRead,
		UpdateContext: resourceStreamUpdate,
		DeleteContext: resourceStreamDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"stream_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"message_stream_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"server_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

type Stream struct {
	ID                string `json:"ID"`
	Name              string `json:"Name"`
	Description       string `json:"Description"`
	MessageStreamType string `json:"MessageStreamType"`
}

func resourceStreamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/message-streams", nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Postmark-Server-Token", d.Get("server_token").(string))
	stream := Stream{
		ID:                d.Get("stream_id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		MessageStreamType: d.Get("message_stream_type").(string),
	}

	body, err := json.Marshal(stream)
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
	err = json.NewDecoder(res.Body).Decode(&stream)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("stream_id").(string))

	return diags
}

func resourceStreamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceStreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
