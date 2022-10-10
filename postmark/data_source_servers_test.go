package postmark

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	postmarkSDK "github.com/keighl/postmark"
)

func TestAccDataSourceServers_basic(t *testing.T) {
	var server postmarkSDK.Server

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostmarkServersConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostmarkServersExists("data.postmark_servers.all", &server),
				),
			},
		},
	})
}

func testAccCheckPostmarkServersExists(n string, server *postmarkSDK.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Data source not Found: %s", n)
		}
		return nil
	}
}

func testAccPostmarkServersConfig() string {
	return `
	data "postmark_servers" "all" {}
  `
}
