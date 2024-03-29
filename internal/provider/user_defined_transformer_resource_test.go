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
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "source", "generate_card_number"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "config.generate_card_number.valid_luhn", "false"),
				),
			},
		},
	})
}

func TestAcc_Ud_Transformer_From_System(t *testing.T) {
	name := acctest.RandString(10)

	config := fmt.Sprintf(`
data "neosync_system_transformer" "gen_cc" {
  source = "generate_card_number"
}

resource "neosync_user_defined_transformer" "test1" {
  name = "%s"
	description = "this is a test"
	source = data.neosync_system_transformer.gen_cc.source
	config = data.neosync_system_transformer.gen_cc.config
}
`, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("neosync_user_defined_transformer.test1", "id"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "name", name),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "description", "this is a test"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "source", "generate_card_number"),
					resource.TestCheckResourceAttr("neosync_user_defined_transformer.test1", "config.generate_card_number.valid_luhn", "true"),
				),
			},
		},
	})
}
