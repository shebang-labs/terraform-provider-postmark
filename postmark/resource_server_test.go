package postmark

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	postmarkSDK "github.com/keighl/postmark"
)

func TestAccServer_basic(t *testing.T) {
	var server postmarkSDK.Server

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a server
			{
				Config: testAccPostmarkServerConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostmarkServerExists("postmark_server.s1", &server),
					testAccCheckPostmarkServerAttributes(&server, &testAccPostmarkServerExpectedAttributes{
						Name:  "Test 1",
						Color: "blue",
					}),
				),
			},
			// Update the server to change the name and color
			{
				Config: testAccPostmarkServerUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostmarkServerExists("postmark_server.s1", &server),
					testAccCheckPostmarkServerAttributes(&server, &testAccPostmarkServerExpectedAttributes{
						Name:  "Test 2",
						Color: "green",
					}),
				),
			},
		},
	})
}

func testAccCheckPostmarkServerExists(n string, server *postmarkSDK.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		serverID := rs.Primary.ID
		if serverID == "" {
			return fmt.Errorf("No server ID is set")
		}
		conn := testAccProvider.Meta().(*postmarkSDK.Client)

		gotServer, err := conn.GetServer(serverID)
		if err != nil {
			return err
		}
		*server = gotServer
		return nil
	}
}

type testAccPostmarkServerExpectedAttributes struct {
	Name  string
	Color string
}

func testAccCheckPostmarkServerAttributes(server *postmarkSDK.Server, want *testAccPostmarkServerExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if server.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", server.Name, want.Name)
		}

		if server.Color != want.Color {
			return fmt.Errorf("got color %q; want %q", server.Color, want.Color)
		}

		return nil
	}
}

func testAccPostmarkServerConfig() string {
	return `
resource "postmark_server" "s1" {
	name             = "Test 1"
	color         = "blue"
}
  `
}

func testAccPostmarkServerUpdateConfig() string {
	return `
resource "postmark_server" "s1" {
  name             = "Test 2"
  color         = "green"
}
  `
}
