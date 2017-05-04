package kafka

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform/helper/schema"

	kafka "github.com/Shopify/sarama"
)

func topic() *schema.Resource {
	return &schema.Resource{
		Create: topicCreate,
		Update: topicCreate,
		Delete: TopicDelete,
		Read:   TopicRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the topic",
			},

			"partitions": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Description: "The number of partitions the topic should have",
			},

			"replication_factor": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Description: "The replication factor of the topic",
			},
			// Data is passed as JSON so that an arbitrary structure is
			// possible, rather than forcing e.g. all values to be strings.
			"config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

type kafkaTopic struct {
	Name              string
	Partitions        int
	ReplicationFactor int
	Config            map[string]string
}

func parseDToTopic(d *schema.ResourceData) *kafkaTopic {
	topicName := d.Get("name").(string)
	numPartitions := d.Get("partitions").(int)
	replicationFactor := d.Get("replication_factor").(int)
	config := d.Get("config").(map[string]interface{})

	m2 := make(map[string]string)
	for key, value := range config {
		switch value := value.(type) {
		case string:
			m2[key] = value
		}
	}
	return &kafkaTopic{
		Name:              topicName,
		Partitions:        numPartitions,
		ReplicationFactor: replicationFactor,
		Config:            m2,
	}
}

func topicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(kafka.Client)
	timeout := int32(1000)

	t := parseDToTopic(d)
	spew.Dump(t)

	err := client.CreateTopic(t.Name, int32(t.Partitions), int16(t.ReplicationFactor), t.Config, timeout)

	if err != nil {
		return fmt.Errorf("error creating topic: %s", err)
	}

	if err != nil {
		return fmt.Errorf("error getting topics: %s", err)
	}

	return nil
}

func TopicDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func TopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(kafka.Client)
	name := d.Id()

	topics, err := client.Topics()
	if err != nil {
		return err
	}
	for _, t := range topics {
		if t == name {
			client.RefreshMetadata(name)
			md := kafka.MetadataRequest{Version: 1, Topics: []string{name}}

			mdr, err := client.Brokers()[0].GetMetadata(&md)
			if err != nil {
				return err
			}
			actualName := mdr.Topics[0].Name
			partitions := len(mdr.Topics[0].Partitions)
			d.Set("name", actualName)
			d.Set("partitions", partitions)

			return nil
		}
	}
	return nil
}
