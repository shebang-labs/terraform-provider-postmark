package postmark

import (
	"context"
	"fmt"
	"time"

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
	serverToken := d.Get("server_token").(string)
	stream := Stream{
		ID:                d.Get("stream_id").(string),
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		MessageStreamType: d.Get("message_stream_type").(string),
	}
	diags, stream := doStreamRequests("POST", "", stream, serverToken)
	if stream.ID == "" {
		return diags
	}
	d.SetId(d.Get("stream_id").(string))

	return diags
}

func resourceStreamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChange("name") || d.HasChange("description") {
		streamId := d.Id()
		serverToken := d.Get("server_token").(string)
		stream := Stream{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
		_, stream = doStreamRequests("PATCH", fmt.Sprintf("/%s", streamId), stream, serverToken)

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceStreamRead(ctx, d, m)
}

func resourceStreamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverToken := d.Get("server_token").(string)
	streamId := d.Id()
	stream := Stream{}
	diags, stream := doStreamRequests("GET", fmt.Sprintf("/%s", streamId), stream, serverToken)

	if stream.ID == "" {
		return diags
	}
	d.SetId(stream.ID)
	if err := d.Set("name", stream.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", stream.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("message_stream_type", stream.MessageStreamType); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// there is no delete endpoint for streams just archive and unarchive
func resourceStreamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
