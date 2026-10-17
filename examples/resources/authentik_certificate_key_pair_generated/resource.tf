# Have authentik generate a certificate-key pair, e.g. to sign OAuth2 tokens
# (`signing_key` on `authentik_provider_oauth2`). The private key stays in
# authentik and never reaches the Terraform state.

resource "authentik_certificate_key_pair_generated" "signing" {
  name          = "oauth2-signing"
  common_name   = "oauth2-signing-2026"
  validity_days = 730
}
