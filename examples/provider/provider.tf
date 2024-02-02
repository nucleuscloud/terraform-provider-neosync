# Configure the Neosync provider using the required_providers stanza.
# You may optionally use a version directive to prevent breaking
# changes occurring unannounced.
terraform {
  required_providers {
    vercel = {
      source  = "neosync/neosync"
      version = "~> 0.3"
    }
  }
}

provider "neosync" {
  # Or omirt this for the endpoint to be read
  # from the NEOSYNC_ENDPOINT environment variable
  endpoint = var.neosync_endpoint

  # Or omit this for the api_token to be read
  # from the NEOSYNC_API_TOKEN environment variable
  # or if running Neosync in unauthenticated mode, omit entirely.
  # If running in unauth mode, the account id must be provided in some fashion
  api_token = var.neosync_api_token

  # Optional account id
  # This can be inferred from the API Key, or if the account_id is provided on the resource
  account_id = var.neosync_account_id
}
