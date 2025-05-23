package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"neosync": providerserver.NewProtocol6WithError(New("test", "http://localhost:8080")()),
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv(apiTokenEnvVarKey) == "" {
		mustHaveEnv(t, accountIdEnvVarKey)
	} else {
		mustHaveEnv(t, apiTokenEnvVarKey)
	}
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func mustHaveEnv(t *testing.T, name string) {
	if os.Getenv(name) == "" {
		t.Fatalf("%s environment variable must be set for acceptance tests", name)
	}
}

// Retrieves the account_id from state during a terraform check. Mutates the input accountId.
func GetAccountIdFromState(resource string, onAccountId func(accountId string)) func(s *terraform.State) error {
	return GetAttributeFromState(resource, "account_id", onAccountId)
}

func GetAttributeFromState(resource string, attribute string, onAttribute func(attribute string)) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		attr := rs.Primary.Attributes[attribute]
		onAttribute(attr)
		return nil
	}
}

// Ensures that the value did not change between runs.
func GetTestAccountIdFromStateFn(resource string, getAccountId func() string) func(s *terraform.State) error {
	return GetTestAttributeFromStateFn(resource, "account_id", getAccountId)
}

func GetTestAttributeFromStateFn(resource string, attributeKey string, getAttribute func() string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		attribute := getAttribute()
		if rs.Primary.Attributes[attributeKey] != attribute {
			return fmt.Errorf("%s changed unexpectedly. Was %s, now %s", attributeKey, attribute, rs.Primary.Attributes[attributeKey])
		}
		return nil
	}
}
