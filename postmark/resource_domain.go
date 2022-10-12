package postmark

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"return_path_domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

type Domain struct {
	Name             string `json:"Name"`
	ID               int64  `json:"ID"`
	ReturnPathDomain string `json:"ReturnPathDomain"`
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	domain := Domain{
		Name:             d.Get("name").(string),
		ReturnPathDomain: d.Get("return_path_domain").(string),
	}
	diags, domain := doDomainRequests("POST", "", domain, m)
	if domain.ID == 0 {
		return diags
	}
	d.SetId(strconv.FormatInt(domain.ID, 10))

	return diags
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	if d.HasChange("return_path_domain") {
		domainId := d.Id()
		domain := Domain{
			ReturnPathDomain: d.Get("return_path_domain").(string),
		}
		_, domain = doDomainRequests("PUT", domainId, domain, m)
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceDomainRead(ctx, d, m)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainId := d.Id()
	domain := Domain{}
	diags, domain := doDomainRequests("GET", domainId, domain, m)

	if domain.ID == 0 {
		return diags
	}
	d.SetId(strconv.FormatInt(domain.ID, 10))
	if err := d.Set("name", domain.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainId := d.Id()
	diags, _ := doDomainRequests("DELETE", domainId, Domain{}, m)
	d.SetId("")
	return diags
}
