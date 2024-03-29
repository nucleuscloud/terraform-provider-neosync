package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Ud_Transformer_DataSource(t *testing.T) {
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

data "neosync_user_defined_transformer" "test1" {
	id = neosync_user_defined_transformer.test1.id
}
`, name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.neosync_user_defined_transformer.test1", "id"),
					resource.TestCheckResourceAttr("data.neosync_user_defined_transformer.test1", "name", name),
				),
			},
		},
	})
}
