package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceProviderOAuth2(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderOAuth2(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "allowed_redirect_uris.0.redirect_uri_type", "authorization"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "allowed_redirect_uris.1.redirect_uri_type", "logout"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName),
				),
			},
			{
				Config: testAccResourceProviderOAuth2(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName+"test"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "allowed_redirect_uris.0.redirect_uri_type", "authorization"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "allowed_redirect_uris.1.redirect_uri_type", "logout"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName+"test"),
				),
			},
		},
	})
}

func TestAccResourceProviderOAuth2_WithSecret(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	appName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderOAuth2WithSecret(rName, appName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName),
				),
			},
			{
				Config: testAccResourceProviderOAuth2WithSecret(rName+"test", appName+"test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "name", rName+"test"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_id", rName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "name", appName+"test"),
					resource.TestCheckResourceAttr("authentik_application.name", "slug", appName+"test"),
				),
			},
		},
	})
}

func testAccResourceProviderOAuth2(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}
data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

data "authentik_certificate_key_pair" "generated" {
  name = "authentik Self-signed Certificate"
  fetch_key = false
  fetch_certificate = false
}
resource "authentik_provider_oauth2" "name" {
  name      = "%[1]s"
  client_id = "%[1]s"
  # client_secret = "test"
  signing_key = data.authentik_certificate_key_pair.generated.id
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
  allowed_redirect_uris = [
    {
      matching_mode     = "strict",
      url               = "http://localhost/callback",
      redirect_uri_type = "authorization",
    },
    {
      matching_mode     = "strict",
      url               = "http://localhost/",
      redirect_uri_type = "logout",
    }
  ]
}

resource "authentik_application" "name" {
  name              = "%[2]s"
  slug              = "%[2]s"
  protocol_provider = authentik_provider_oauth2.name.id
}
`, name, appName)
}

func testAccResourceProviderOAuth2WithSecret(name string, appName string) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}
data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

data "authentik_certificate_key_pair" "generated" {
  name = "authentik Self-signed Certificate"
  fetch_key = false
  fetch_certificate = false
}
resource "authentik_provider_oauth2" "name" {
  name      = "%[1]s"
  client_id = "%[1]s"
  client_secret = "test"
  signing_key = data.authentik_certificate_key_pair.generated.id
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow = data.authentik_flow.default-provider-invalidation-flow.id
}

resource "authentik_application" "name" {
  name              = "%[2]s"
  slug              = "%[2]s"
  protocol_provider = authentik_provider_oauth2.name.id
}
`, name, appName)
}

func TestAccResourceProviderOAuth2_WithSecretWriteOnly(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderOAuth2WithSecretWriteOnly(rName, "first", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_secret", ""),
					resource.TestCheckNoResourceAttr("authentik_provider_oauth2.name", "client_secret_wo"),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_secret_wo_version", "1"),
					testAccCheckProviderOAuth2ClientSecret("authentik_provider_oauth2.name", "first"),
				),
			},
			{
				// A new value without a new version is not sent.
				Config: testAccResourceProviderOAuth2WithSecretWriteOnly(rName, "second", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_secret", ""),
					testAccCheckProviderOAuth2ClientSecret("authentik_provider_oauth2.name", "first"),
				),
			},
			{
				Config: testAccResourceProviderOAuth2WithSecretWriteOnly(rName, "second", 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_secret", ""),
					resource.TestCheckResourceAttr("authentik_provider_oauth2.name", "client_secret_wo_version", "2"),
					testAccCheckProviderOAuth2ClientSecret("authentik_provider_oauth2.name", "second"),
				),
			},
		},
	})
}

// Reads the secret authentik actually holds, since state no longer has it.
func testAccCheckProviderOAuth2ClientSecret(name string, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}
		p := Provider("test", false)
		ac, diags := p.ConfigureContextFunc(context.Background(), schema.TestResourceDataRaw(&testing.T{}, p.Schema, map[string]any{
			"url":   os.Getenv("AUTHENTIK_URL"),
			"token": os.Getenv("AUTHENTIK_TOKEN"),
		}))
		if diags.HasError() {
			return fmt.Errorf("configuring client: %v", diags)
		}
		res, _, err := ac.(*APIClient).client.ProvidersAPI.ProvidersOauth2Retrieve(context.Background(), int32(id)).Execute()
		if err != nil {
			return err
		}
		if res.GetClientSecret() != want {
			return fmt.Errorf("client secret in authentik is %q, want %q", res.GetClientSecret(), want)
		}
		return nil
	}
}

func testAccResourceProviderOAuth2WithSecretWriteOnly(name string, secret string, version int) string {
	return fmt.Sprintf(`
data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}
data "authentik_flow" "default-provider-invalidation-flow" {
  slug = "default-provider-invalidation-flow"
}

resource "authentik_provider_oauth2" "name" {
  name                     = "%[1]s"
  client_id                = "%[1]s"
  client_secret_wo         = "%[2]s"
  client_secret_wo_version = %[3]d
  authorization_flow       = data.authentik_flow.default-authorization-flow.id
  invalidation_flow        = data.authentik_flow.default-provider-invalidation-flow.id
}
`, name, secret, version)
}
