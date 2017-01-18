package github

import (
	"log"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/schema"
)

func resouceGithubOrganizationWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubOrganizationWebhookCreate,
		Read:   resourceGithubOrganizationWebhookRead,
		Update: resourceGithubOrganizationWebhookUpdate,
		Delete: resourceGithubOrganizationWebhookDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
			},
			"events": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceGithubOrganizationWebhookObject(d *schema.ResourceData) *github.Hook {
	name := d.Get("name").(string)
	//events := d.Get("events").([]interface{})
	active := d.Get("active").(bool)
	config := d.Get("config").(map[string]interface{})
	d.GetOk("active")

	//var moo []string
	//for _, e := range events {
	//moo = append(moo, e.(string))
	//}
	//if len(moo) == 0 {
	//moo = []string{"push"}
	//}

	repo := &github.Hook{
		Name:   &name,
		Config: config,
		//Events: moo,
		//Active: &active,
	}

	return repo
}

func resourceGithubOrganizationWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Organization).client

	hook, r, err := client.Organizations.CreateHook(meta.(*Organization).name, resourceGithubOrganizationWebhookObject(d))

	//defer resp.Body.Close()

	log.Printf("[DEBUG] %s", r.StatusCode)

	if err != nil {
		log.Printf("[DEBUG] could not create github organisation hook for %s -> %s ", meta.(*Organization).name, err)
		return err
	}

	d.SetId(strconv.Itoa(*hook.ID))

	return resourceGithubOrganizationWebhookRead(d, meta)
}

func resourceGithubOrganizationWebhookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Organization).client
	hookID, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Printf("[ERROR] Could not convert %s to int: %s", d.Id(), err)
		return err
	}

	log.Printf("[DEBUG] read webhook %d for github org %s", hookID, meta.(*Organization).name)

	hook, resp, err := client.Organizations.GetHook(meta.(*Organization).name, hookID)

	if err != nil {
		if resp.StatusCode == 404 {
			log.Printf(
				"[WARN] removing %s organisation webhook, it no longer exists in github",
				meta.(*Organization).name,
			)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", hook.Name)
	d.Set("config", hook.Config)
	d.Set("events", hook.Events)
	d.Set("active", hook.Active)

	return nil
}

func resourceGithubOrganizationWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Organization).client
	hookReq := resourceGithubOrganizationWebhookObject(d)
	hookID, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Printf("[ERROR] Could not convert %s to int: %s", d.Id(), err)
		return err
	}

	hook, _, err := client.Organizations.EditHook(meta.(*Organization).name, hookID, hookReq)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(*hook.ID))

	return resourceGithubOrganizationWebhookRead(d, meta)
}

func resourceGithubOrganizationWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Organization).client
	hookID, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Printf("[ERROR] Could not convert %s to int: %s", d.Id(), err)
		return err
	}

	_, err = client.Organizations.DeleteHook(meta.(*Organization).name, hookID)
	return err
}
