package typesense

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testClusterId = os.Getenv("CLUSTER_ID")

func TestClusterApiKeysResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "typesense_cluster_api_keys" "test" {
	cluster_id = "%s"
}
`, testClusterId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_cluster_api_keys.test", "cluster_id", testClusterId),
					resource.TestCheckResourceAttrSet("typesense_cluster_api_keys.test", "admin_key"),
					resource.TestCheckResourceAttrSet("typesense_cluster_api_keys.test", "search_only_key"),
				),
			},
			// ImportState testing
			// Update and Read testing
			// Delete testing automatically occurs in TestCase
		},
	})
}
