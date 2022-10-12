package postmark

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	postmarkSDK "github.com/keighl/postmark"
)

func TestAccDomain_basic(t *testing.T) {
	var domain Domain

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a domain
			{
				Config: testAccPostmarkDomainConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostmarkDomainExists("postmark_domain.d1", &domain),
					testAccCheckPostmarkDomainAttributes(&domain, &testAccPostmarkDomainExpectedAttributes{
						Name:             "rentola.test",
						ReturnPathDomain: "testoooo.rentola.test",
					}),
				),
			},
			// Update the domain to change the return path domain
			{
				Config: testAccPostmarkDomainUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostmarkDomainExists("postmark_domain.d1", &domain),
					testAccCheckPostmarkDomainAttributes(&domain, &testAccPostmarkDomainExpectedAttributes{
						Name:             "rentola.test",
						ReturnPathDomain: "test.rentola.test",
					}),
				),
			},
		},
	})
}

func testAccCheckPostmarkDomainExists(n string, domain *Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		domainID := rs.Primary.ID
		if domainID == "" {
			return fmt.Errorf("No domain ID is set")
		}

		diags, gotDomain := doDomainRequests("GET", domainID, Domain{}, testAccProvider.Meta().(*postmarkSDK.Client))
		if diags != nil {
			return fmt.Errorf("Cann't get the domain")
		}
		*domain = gotDomain
		return nil
	}
}

type testAccPostmarkDomainExpectedAttributes struct {
	Name             string
	ReturnPathDomain string
}

func testAccCheckPostmarkDomainAttributes(domain *Domain, want *testAccPostmarkDomainExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if domain.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", domain.Name, want.Name)
		}

		if domain.ReturnPathDomain != want.ReturnPathDomain {
			return fmt.Errorf("got return path domain %q; want %q", domain.ReturnPathDomain, want.ReturnPathDomain)
		}

		return nil
	}
}

func testAccPostmarkDomainConfig() string {
	return `
resource "postmark_domain" "d1" {
	name = "rentola.test"
  return_path_domain = "testoooo.rentola.test"
}
  `
}

func testAccPostmarkDomainUpdateConfig() string {
	return `
resource "postmark_domain" "d1" {
	name = "rentola.test"
  return_path_domain = "test.rentola.test"
}
  `
}
