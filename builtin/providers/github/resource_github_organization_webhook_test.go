package github

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/google/go-github/github"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGithubOrganizationWebhook_basic(t *testing.T) {
	var hook github.Hook

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGithubOrganizationWebhookDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGithubOrganizationWebhookConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubOrganizationWebhookExists(&hook),
				),
			},
			resource.TestStep{
				Config: testAccGithubOrganizationWebhookUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGithubOrganizationWebhookExists("github_organization_webhook.foo", &repo),
					testAccCheckGithubOrganizationWebhookAttributes(&hook, &testAccGithubOrganizationWebhookExpectedAttributes{
						Name:          "foo",
						Description:   "Terraform acceptance tests!",
						Homepage:      "http://example.com/",
						DefaultBranch: "master",
					}),
				),
			},
		},
	})
}

func TestAccGithubOrganizationWebhook_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGithubOrganizationWebhookDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGithubOrganizationWebhookConfig,
			},
			resource.TestStep{
				ResourceName:      "github_repository.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGithubOrganizationWebhookExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		repoName := rs.Primary.ID
		if repoName == "" {
			return fmt.Errorf("No repository name is set")
		}

		org := testAccProvider.Meta().(*Organization)
		conn := org.client
		gotRepo, _, err := conn.Repositories.Get(org.name, repoName)
		if err != nil {
			return err
		}
		*repo = *gotRepo
		return nil
	}
}

type testAccGithubOrganizationWebhookExpectedAttributes struct {
	Name         string
	Description  string
	Homepage     string
	Private      bool
	HasIssues    bool
	HasWiki      bool
	HasDownloads bool

	DefaultBranch string
}

func testAccCheckGithubOrganizationWebhookAttributes(repo *github.Repository, want *testAccGithubOrganizationWebhookExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if *repo.Name != want.Name {
			return fmt.Errorf("got repo %q; want %q", *repo.Name, want.Name)
		}
		if *repo.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", *repo.Description, want.Description)
		}
		if *repo.Homepage != want.Homepage {
			return fmt.Errorf("got homepage URL %q; want %q", *repo.Homepage, want.Homepage)
		}
		if *repo.Private != want.Private {
			return fmt.Errorf("got private %#v; want %#v", *repo.Private, want.Private)
		}
		if *repo.HasIssues != want.HasIssues {
			return fmt.Errorf("got has issues %#v; want %#v", *repo.HasIssues, want.HasIssues)
		}
		if *repo.HasWiki != want.HasWiki {
			return fmt.Errorf("got has wiki %#v; want %#v", *repo.HasWiki, want.HasWiki)
		}
		if *repo.HasDownloads != want.HasDownloads {
			return fmt.Errorf("got has downloads %#v; want %#v", *repo.HasDownloads, want.HasDownloads)
		}

		if *repo.DefaultBranch != want.DefaultBranch {
			return fmt.Errorf("got default branch %q; want %q", *repo.DefaultBranch, want.DefaultBranch)
		}

		return nil
	}
}

func testAccCheckGithubOrganizationWebhookDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*Organization).client
	orgName := testAccProvider.Meta().(*Organization).name
	hookID, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Printf("[ERROR] Could not convert %s to int: %s", d.Id(), err)
		return err
	}

	for _, rs := range s.RootModule().Resources {
		gotHook, resp, err := conn.Organizations.GetHook(orgName, hookID)

		if err == nil {
			if gotRepo != nil && *gotRepo.Name == rs.Primary.ID {
				return fmt.Errorf("Repository still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

const testAccGithubOrganizationWebhookConfig = `
resource "github_organization_webhook" "foo" {
  name = "compliance_webhook"
	url = "http://mongey.net"
}
`

const testAccGithubOrganizationWebhookUpdateConfig = `
resource "github_organization_webhook" "foo" {
	name = "new_webhook"
	url = "http://mongey.net"
}
`
