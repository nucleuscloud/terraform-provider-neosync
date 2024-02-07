package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Ud_Transformer(t *testing.T) {
	name := acctest.RandString(10)

	testAccConnectionConfig := fmt.Sprintf(`
resource "neosync_user_defined_transformer" "test1" {
  name = "%s"
	description = "this is a test"
	datatype = "int64"
	source = "generate_card_number"
	config = {
		"generate_card_number" = {
			valid_luhn = true
		}
	}
}
`, name)
	testAccConnectionConfigUpdated := fmt.Sprintf(`
resource "neosync_user_defined_transformer" "test1" {
  name = "%s"
	description = "this is a test2"
	datatype = "int64"
	source = "generate_card_number"
	config = {
		"generate_card_number" = {
			valid_luhn = false
		}
	}
}
`, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_user_defined_transformer.test1", "id"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "name", name),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "description", "this is a test"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "datatype", "int64"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "source", "generate_card_number"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "config.generate_card_number.valid_luhn", "true"),
				),
			},
			{
				Config: testAccConnectionConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_user_defined_transformer.test1", "id"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "name", name),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "description", "this is a test2"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "datatype", "int64"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "source", "generate_card_number"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "config.generate_card_number.valid_luhn", "false"),
				),
			},
		},
	})
}
