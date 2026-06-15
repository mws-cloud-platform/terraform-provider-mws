package acctest

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"go.mws.cloud/util-toolset/pkg/utils/consterr"
)

const (
	// ErrResourceNotFound is returned when a terraform resource is not found.
	ErrResourceNotFound = consterr.Error("resource not found")
)

// ResourceCheck is a function that checks the state of a resource by id.
type ResourceCheck func(ctx context.Context, id string) error

// SingleResourceTestCase represents an acceptance test case for a single
// resource to verify that a single resource can be correctly created, updated
// and destroyed. Specifically, it checks that:
//   - resource is planned for creation
//   - resource exists in the state after apply
//   - resource is really created during apply
//   - no changes are planned for the same config on the next apply
//   - data source correctly gets information about the created resource
//   - resource is planned for update if the update configuration is provided
//   - resource is really exists after update
//   - resource is really deleted during destroy
//
// If provider configuration is not empty, it must be valid Terraform provider
// configuration for the resource being tested.
//
// Resource configuration must be valid Terraform configuration with a single
// resource definition. For example:
//
//	resource "mws_compute_virtual_machine" "vm" {
//	  # ...
//	}
//
// DataSource configuration must be valid Terraform configuration with a single
// datasource definition. For example:
//
//	data "mws_compute_virtual_machine" "vm_data" {
//	  # ...
//	}
//
// Update configuration must be valid Terraform configuration with the same
// resource name as the original configuration.
type SingleResourceTestCase struct {
	ProviderConfig        string        // Terraform provider configuration used in each test case step
	ResourceConfig        string        // Terraform configuration for the resource to be tested
	UpdatedResourceConfig string        // Terraform configuration for the updated resource (optional)
	DataSourceConfig      string        // Terraform configuration for the data source to be tested (optional)
	RecreateOnUpdate      bool          // Indicates if resource should be destroyed and re-created on update
	ResourceExists        ResourceCheck // Function to check if resource exists
	ResourceNotExists     ResourceCheck // Function to check if resource does not exist
}

// Build builds a [resource.TestCase] with the resource verification test steps.
func (tc SingleResourceTestCase) Build(ctx context.Context) (resource.TestCase, error) {
	resourceName, err := parseResourceName(tc.ResourceConfig)
	if err != nil {
		return resource.TestCase{}, fmt.Errorf("parse resource name from config: %w", err)
	}

	resourceExists := func(s *terraform.State) error {
		r, err := getResourceFromState(s, resourceName)
		if err != nil {
			return err
		}

		return tc.ResourceExists(ctx, r.Primary.ID)
	}

	steps := []resource.TestStep{
		{
			Config: tc.ResourceConfig,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
				},
			},
			Check: resourceExists,
		},
		{
			Config: tc.ResourceConfig + "\n" + tc.DataSourceConfig,
			// verify that no changes are planned for the same config
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				},
			},
			Check: func(state *terraform.State) error {
				if tc.DataSourceConfig == "" {
					return nil
				}

				dsName, err := parseDataSourceName(tc.DataSourceConfig)
				if err != nil {
					return err
				}

				resName, err := parseResourceName(tc.ResourceConfig)
				if err != nil {
					return err
				}

				return resource.TestCheckResourceAttrPair(
					dsName, "metadata.id",
					resName, "metadata.id",
				)(state)
			},
		},
	}

	if tc.UpdatedResourceConfig != "" {
		updateConfigResourceName, err := parseResourceName(tc.UpdatedResourceConfig)
		if err != nil {
			return resource.TestCase{}, fmt.Errorf("parse resource name from update config: %w", err)
		}
		if updateConfigResourceName != resourceName {
			return resource.TestCase{}, fmt.Errorf("update config has different resource name %q: %w",
				updateConfigResourceName, ErrTerraformInvalidConfig)
		}

		expectedResourceAction := plancheck.ResourceActionUpdate
		if tc.RecreateOnUpdate {
			expectedResourceAction = plancheck.ResourceActionDestroyBeforeCreate
		}

		steps = append(steps, resource.TestStep{
			Config: tc.UpdatedResourceConfig,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, expectedResourceAction),
				},
			},
			Check: resourceExists,
		})
	}

	if tc.ProviderConfig != "" {
		for i := range steps {
			steps[i].Config = tc.ProviderConfig + "\n" + steps[i].Config
		}
	}

	checkDestroy := func(s *terraform.State) error {
		r, err := getResourceFromState(s, resourceName)
		if err != nil {
			return err
		}
		return tc.ResourceNotExists(ctx, r.Primary.ID)
	}

	return resource.TestCase{
		Steps:        steps,
		CheckDestroy: checkDestroy,
	}, nil
}

func getResourceFromState(s *terraform.State, resourceName string) (*terraform.ResourceState, error) {
	r, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("get %q from state: %w", resourceName, ErrResourceNotFound)
	}
	return r, nil
}

func parseResourceName(config string) (string, error) {
	return parseBlockName(config, "resource")
}

func parseDataSourceName(config string) (string, error) {
	name, err := parseBlockName(config, "data")
	if err != nil {
		return "", err
	}
	return "data." + name, nil
}

func parseBlockName(config, blockType string) (string, error) {
	var name string
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCL([]byte(config), "config")
	if diags.HasErrors() {
		return "", fmt.Errorf("parse config: %w", diags)
	}

	body := f.Body.(*hclsyntax.Body)
	for _, block := range body.Blocks {
		if block.Type != blockType {
			continue
		}
		if len(block.Labels) != 2 {
			return "", fmt.Errorf("invalid block %s: %w", block.Range(), ErrTerraformInvalidConfig)
		}
		if name != "" {
			return "", fmt.Errorf("multiple blocks found: %w", ErrTerraformInvalidConfig)
		}
		name = block.Labels[0] + "." + block.Labels[1]
	}
	if name == "" {
		return "", fmt.Errorf("no block found: %w", ErrTerraformInvalidConfig)
	}

	return name, nil
}
