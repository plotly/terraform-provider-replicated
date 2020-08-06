package replicated

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLicense() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLicenseRead,
		Schema: map[string]*schema.Schema{
			"customer_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"license_base64": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLicenseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiToken := m.(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	customerID := d.Get("customer_id").(string)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.replicated.com/vendor/v1/licensekey/%s", customerID), nil)
	req.Header.Set("Authorization", apiToken)

	if err != nil {
		return diag.FromErr(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("license_base64", string(body)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(customerID)

	return diags
}
