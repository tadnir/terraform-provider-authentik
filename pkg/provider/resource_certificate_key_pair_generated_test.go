package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceCertificateKeyPairGenerated(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	var firstID string
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCertificateKeyPairGenerated(rName, rName+"-cn"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_certificate_key_pair_generated.name", "name", rName),
					resource.TestMatchResourceAttr("authentik_certificate_key_pair_generated.name", "certificate_data", regexp.MustCompile("BEGIN CERTIFICATE")),
					resource.TestCheckResourceAttrSet("authentik_certificate_key_pair_generated.name", "fingerprint_sha256"),
					resource.TestCheckNoResourceAttr("authentik_certificate_key_pair_generated.name", "key_data"),
					func(s *terraform.State) error {
						firstID = s.RootModule().Resources["authentik_certificate_key_pair_generated.name"].Primary.ID
						return nil
					},
				),
			},
			// Renaming keeps the key.
			{
				Config: testAccResourceCertificateKeyPairGenerated(rName+"test", rName+"-cn"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_certificate_key_pair_generated.name", "name", rName+"test"),
					func(s *terraform.State) error {
						if id := s.RootModule().Resources["authentik_certificate_key_pair_generated.name"].Primary.ID; id != firstID {
							return fmt.Errorf("rename replaced the key pair: %s != %s", id, firstID)
						}
						return nil
					},
				),
			},
			// A new common name generates a new key.
			{
				Config: testAccResourceCertificateKeyPairGenerated(rName+"test", rName+"-cn2"),
				Check: func(s *terraform.State) error {
					if id := s.RootModule().Resources["authentik_certificate_key_pair_generated.name"].Primary.ID; id == firstID {
						return fmt.Errorf("common_name change kept the old key pair")
					}
					return nil
				},
			},
		},
	})
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerTestFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceCertificateKeyPairGenerated(rName, rName+"-cn"),
				ExpectError: regexp.MustCompile("mock-failed-request"),
			},
		},
	})
}

func testAccResourceCertificateKeyPairGenerated(name string, cn string) string {
	return fmt.Sprintf(`
resource "authentik_certificate_key_pair_generated" "name" {
  name          = "%[1]s"
  common_name   = "%[2]s"
  validity_days = 30
  alg           = "ecdsa"
}
`, name, cn)
}
