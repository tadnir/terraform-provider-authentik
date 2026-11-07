package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceProviderSCIM(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderSCIM(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_scim.name", "name", rName),
				),
			},
			{
				Config: testAccResourceProviderSCIM(rName + "test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_provider_scim.name", "name", rName+"test"),
				),
			},
		},
	})
}

func testAccResourceProviderSCIM(name string) string {
	return fmt.Sprintf(`
resource "authentik_provider_scim" "name" {
  name  = "%[1]s"
  url   = "http://localhost"
  token = "foo"
}
`, name)
}

func TestAccResourceProviderSCIM_WithTokenWriteOnly(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProviderSCIMWithTokenWriteOnly(rName, "first", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("authentik_provider_scim.name", "token"),
					resource.TestCheckNoResourceAttr("authentik_provider_scim.name", "token_wo"),
					resource.TestCheckResourceAttr("authentik_provider_scim.name", "token_wo_version", "1"),
				),
			},
			{
				// A new value without a new version is not sent, and plans no change.
				Config:   testAccResourceProviderSCIMWithTokenWriteOnly(rName, "second", 1),
				PlanOnly: true,
			},
			{
				Config: testAccResourceProviderSCIMWithTokenWriteOnly(rName, "second", 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("authentik_provider_scim.name", "token"),
					resource.TestCheckResourceAttr("authentik_provider_scim.name", "token_wo_version", "2"),
				),
			},
		},
	})
}

func testAccResourceProviderSCIMWithTokenWriteOnly(name string, token string, version int) string {
	return fmt.Sprintf(`
resource "authentik_provider_scim" "name" {
  name             = "%[1]s"
  url              = "http://localhost"
  token_wo         = "%[2]s"
  token_wo_version = %[3]d
}
`, name, token, version)
}
