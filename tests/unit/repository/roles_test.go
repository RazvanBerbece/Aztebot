package repository_test

import (
	"testing"

	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/internal/bot-service/data/repositories"
)

func TestGetRoleById(t *testing.T) {

	// Arrange
	cases := []struct {
		input          int
		expectedOutput dataModels.Role
	}{
		{1, dataModels.Role{
			Id:          1,
			RoleName:    "aztec",
			DisplayName: "Aztec",
		}},
		{11, dataModels.Role{
			Id:          11,
			RoleName:    "practicus",
			DisplayName: "🎩 Practicus",
		}},
		{18, dataModels.Role{
			Id:          18,
			RoleName:    "ipsissimus",
			DisplayName: "⚔️ Ipsissimus",
		}},
		{19, dataModels.Role{
			Id:          19,
			RoleName:    "arhitect",
			DisplayName: "👁‍🗨 Arhitect",
		}},
	}

	repoToTest := repositories.NewRolesRepository()

	// Act & Assert
	for _, c := range cases {
		// Act
		output, err := repoToTest.GetRoleById(c.input)
		// Assert
		if err != nil && c.input > 0 {
			t.Errorf("Role expected, error occurred instead: %v", err)
		}
		if output.Id != c.expectedOutput.Id {
			t.Errorf("incorrect output for `Id`: expected `%d` but got `%d`", c.expectedOutput.Id, output.Id)
		}
		if output.RoleName != c.expectedOutput.RoleName {
			t.Errorf("incorrect output for `RoleName`: expected `%s` but got `%s`", c.expectedOutput.RoleName, output.RoleName)
		}
		if output.DisplayName != c.expectedOutput.DisplayName {
			t.Errorf("incorrect output for `DisplayName`: expected `%s` but got `%s`", c.expectedOutput.DisplayName, output.DisplayName)
		}
	}

}

func TestGetRole(t *testing.T) {

	// Arrange
	cases := []struct {
		input          string
		expectedOutput dataModels.Role
	}{
		{"Aztec", dataModels.Role{
			Id:          1,
			RoleName:    "aztec",
			DisplayName: "Aztec",
		}},
		{"👁‍🗨 Arhitect", dataModels.Role{
			Id:          19,
			RoleName:    "arhitect",
			DisplayName: "👁‍🗨 Arhitect",
		}},
		{"Dominus", dataModels.Role{
			Id:          8,
			RoleName:    "dominus",
			DisplayName: "Dominus",
		}},
	}

	repoToTest := repositories.NewRolesRepository()

	// Act & Assert
	for _, c := range cases {
		// Act
		output, err := repoToTest.GetRole(c.input)
		// Assert
		if err != nil {
			t.Errorf("Role expected, error occurred instead: %v", err)
		}
		if output.Id != c.expectedOutput.Id {
			t.Errorf("incorrect output for `Id`: expected `%d` but got `%d`", c.expectedOutput.Id, output.Id)
		}
		if output.RoleName != c.expectedOutput.RoleName {
			t.Errorf("incorrect output for `RoleName`: expected `%s` but got `%s`", c.expectedOutput.RoleName, output.RoleName)
		}
		if output.DisplayName != c.expectedOutput.DisplayName {
			t.Errorf("incorrect output for `DisplayName`: expected `%s` but got `%s`", c.expectedOutput.DisplayName, output.DisplayName)
		}
	}

}
