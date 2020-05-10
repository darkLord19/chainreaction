package validators

import (
	"fmt"

	"github.com/chainreaction/models"
)

// ValidateInstance validates game instance member created after new game req
func ValidateInstance(inst *models.Instance) error {
	if inst.PlayersCount < 2 {
		return fmt.Errorf("At least two players needed")
	}
	if inst.Dimension == 0 {
		return fmt.Errorf("Invalid dimension value: %v", inst.Dimension)
	}
	return nil
}
