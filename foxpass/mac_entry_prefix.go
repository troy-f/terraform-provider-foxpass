package foxpass

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func MacEntryPrefix() *schema.Resource {
	return &schema.Resource{
		CreateContext: macEntryPrefixCreate,
		ReadContext:   macEntryPrefixRead,
		DeleteContext: macEntryPrefixDelete,
		Schema: map[string]*schema.Schema{
			"entryname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func macEntryPrefixCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*FoxpassClient)

	var diags diag.Diagnostics

	entryname := d.Get("entryname").(string)
	prefix := d.Get("prefix").(string)

	err := c.AddMacEntryPrefix(entryname, prefix)
	if err != nil {
		return diag.FromErr(err)
	}

	exists, err := c.GetMacEntryPrefix(entryname, prefix)
	if err != nil {
		return diag.FromErr(err)
	}
	if !exists {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create mac entry prefix",
			Detail:   "Unable to create mac entry prefix",
		})
		return diags
	}

	d.SetId(entryname + "-" + prefix)

	return diags
}

func macEntryPrefixRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*FoxpassClient)

	var diags diag.Diagnostics
	
	entryname := d.Get("entryname").(string)
	prefix := d.Get("prefix").(string)

	exists, err := c.GetMacEntryPrefix(entryname, prefix)
	if err != nil {
		return diag.FromErr(err)
	}

	if !exists {
		d.Set("entryname", "")
		d.Set("prefix", "")
		d.SetId("")
	}

	return diags
}

func macEntryPrefixDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*FoxpassClient)

	var diags diag.Diagnostics
	
	entryname := d.Get("entryname").(string)
	prefix := d.Get("prefix").(string)

	err := c.DeleteMacEntryPrefix(entryname, prefix)
	if err != nil {
		return diag.FromErr(err)
	}

	exists, err := c.GetMacEntryPrefix(entryname, prefix)
	if err != nil {
		return diag.FromErr(err)
	}
	if exists {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete mac entry prefix",
			Detail:   "Unable to delete mac entry prefix",
		})
		return diags
	}

    d.SetId("")

	return diags
}
