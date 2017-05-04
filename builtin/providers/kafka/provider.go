package kafka

import (
	"fmt"
	"time"

	kafka "github.com/Shopify/sarama"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				//TODO: Make a list
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_BOOTSTRAP_SERVER_ADDR", nil),
				Description: "URL of the root of the target Vault server.",
			},
		},

		ConfigureFunc: providerConfigure,

		ResourcesMap: map[string]*schema.Resource{
			"kafka_topic": topic(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	address := d.Get("address").(string)

	config := kafka.NewConfig()
	config.Net.DialTimeout = 30 * time.Second
	config.Version = kafka.V0_10_2_1
	client, err := kafka.NewClient([]string{address}, config)

	if err != nil {
		return nil, fmt.Errorf("failed to configure Kafka Client: %s", err)
	}

	return client, nil
}
