package kafka

import (
	"fmt"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestResourceGenericSecret(t *testing.T) {
	r.Test(t, r.TestCase{
		Providers: testProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []r.TestStep{
			r.TestStep{
				Config: testTopic_initialConfig,
				Check:  testTopic_initialCheck,
			},
			r.TestStep{
				Config: testTopic_updateConfig,
				Check:  testTopic_updateCheck,
			},
		},
	})
}

var testTopic_initialConfig = `
resource "kafka_topic" "logs" {
    name = "abc1234"
		partitions = 1
		replication_factor = 1
}
`

func testTopic_initialCheck(s *terraform.State) error {
	resourceState := s.Modules[0].Resources["kafka_topic.logs"]
	if resourceState == nil {
		return fmt.Errorf("resource not found in state")
	}

	instanceState := resourceState.Primary
	if instanceState == nil {
		return fmt.Errorf("resource has no primary instance")
	}

	path := instanceState.ID

	if path != instanceState.Attributes["name"] {
		return fmt.Errorf("name doesn't match path")
	}
	if path != "abc1234" {
		return fmt.Errorf("unexpected topic")
	}

	return nil
}

var testTopic_updateConfig = `
resource "kafka_topic" "logs" {
    name = "abc1234"
		partitions = 2
		replication_factor = 1
}
`

func testTopic_updateCheck(s *terraform.State) error {

	return nil
}
