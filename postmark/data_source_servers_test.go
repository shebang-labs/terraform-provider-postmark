package postmark

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceServers_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostmarkServersConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostmarkServersExists("data.postmark_servers.all"),
				),
			},
		},
	})
}

func testAccCheckPostmarkServersExists(n string) resource.TestCheckFunc {
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
