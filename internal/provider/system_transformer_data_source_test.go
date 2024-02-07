package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_System_Transformer_DataSource(t *testing.T) {
	config := `
data "neosync_system_transformer" "gen_cc" {
  source = "generate_card_number"
}
	`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.neosync_system_transformer.gen_cc", "source"),
					resource.TestCheckResourceAttrSet("data.neosync_system_transformer.gen_cc", "name"),
					resource.TestCheckResourceAttrSet("data.neosync_system_transformer.gen_cc", "description"),
					resource.TestCheckResourceAttrSet("data.neosync_system_transformer.gen_cc", "datatype"),
					resource.TestCheckResourceAttrSet("data.neosync_system_transformer.gen_cc", "config.generate_card_number.valid_luhn"),
				),
			},
		},
	})
}
