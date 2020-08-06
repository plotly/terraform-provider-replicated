provider "replicated" {
  api_token = "API_TOKEN"
}

data "replicated_license" "customer" {
  customer_id = "CUSTOMER_ID_FROM_URL"
}

output "license_base64" {
  value = data.replicated_license.customer.license_base64
}
