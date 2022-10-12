package postmark

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStream_basic(t *testing.T) {
	var stream Stream

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a stream
			{
				Config: testAccPostmarkStreamConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostmarkStreamExists("postmark_stream.st1", &stream),
					testAccCheckPostmarkStreamAttributes(&stream, &testAccPostmarkStreamExpectedAttributes{
						ID:                "transactional-dev-11",
						Name:              "Stream 1",
						Description:       "This is my first transactional stream",
						MessageStreamType: "Transactional",
					}),
				),
			},
			// Update thes stream to change the name and description
			{
				Config: testAccPostmarkStreamUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostmarkStreamExists("postmark_stream.st1", &stream),
					testAccCheckPostmarkStreamAttributes(&stream, &testAccPostmarkStreamExpectedAttributes{
						ID:                "transactional-dev-11",
						Name:              "Stream 2",
						Description:       "This is my first transactional stream after update",
						MessageStreamType: "Transactional",
					}),
				),
			},
		},
	})
}

func testAccCheckPostmarkStreamExists(n string, stream *Stream) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		streamID := rs.Primary.ID
		if streamID == "" {
			return fmt.Errorf("No stream ID is set")
		}

		diags, gotStream := doStreamRequests("GET", streamID, Stream{}, rs.Primary.Attributes["server_token"])
		if diags != nil {
			return fmt.Errorf("Cann't get the domain")
		}
		*stream = gotStream
		return nil
	}
}

type testAccPostmarkStreamExpectedAttributes struct {
	ID                string
	Name              string
	Description       string
	MessageStreamType string
}

func testAccCheckPostmarkStreamAttributes(stream *Stream, want *testAccPostmarkStreamExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if stream.ID != want.ID {
			return fmt.Errorf("got id %q; want %q", stream.ID, want.ID)
		}

		if stream.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", stream.Name, want.Name)
		}

		if stream.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", stream.Description, want.Description)
		}

		if stream.MessageStreamType != want.MessageStreamType {
			return fmt.Errorf("got message stream type %q; want %q", stream.MessageStreamType, want.MessageStreamType)
		}

		return nil
	}
}

func testAccPostmarkStreamConfig() string {
	return `
resource "postmark_stream" "st1" {
  stream_id = "transactional-dev-11"
  name = "Stream 1"
  description = "This is my first transactional stream"
  message_stream_type = "Transactional"
  server_token = "Server token"
}	  
  `
}

func testAccPostmarkStreamUpdateConfig() string {
	return `
resource "postmark_stream" "st1" {
  stream_id = "transactional-dev-11"
  name = "Stream 2"
  description = "This is my first transactional stream after update"
  message_stream_type = "Transactional"
  server_token = "Server token"
}	  
  `
}
