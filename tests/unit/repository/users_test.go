package repository_test

import (
	"testing"

	repositories "github.com/RazvanBerbece/Aztebot/internal/data/repositories/aztebot"
)

func TestUserExists(t *testing.T) {

	// Arrange
	cases := []struct {
		input          string
		expectedOutput int
	}{
		{"", -1},
		{"ThisIsAFakeUserId", -1},
		{"573659533361020941", 1},
	}

	repoToTest := repositories.NewUsersRepository()

	// Act & Assert
	for _, c := range cases {
		// Act
		output := repoToTest.UserExists(c.input)
		// Assert
		switch c.input {
		case "":
			if output > 0 {
				t.Errorf("Expected to not find user %s in DB", c.input)
			}
		case "573659533361020941":
			if output < 0 {
				t.Errorf("Expected to find user in DB")
			}
		default:
			if output > 0 {
				t.Errorf("Expected to not find user %s in DB", c.input)
			}
		}
	}

}
