package repository_test

// Uncomment this when DB migrations are working in CI or when there is a DB provisioner image available
// func TestGetUserTimeout(t *testing.T) {

// 	// Arrange
// 	repoSource := repositories.NewTimeoutsRepository()

// 	idUserWithActiveTimeout := "1234567890123"
// 	idUserWithNoActiveTimeout := "1234567890"

// 	arrangedTimeout := dataModels.Timeout{
// 		UserId:            idUserWithActiveTimeout,
// 		Reason:            gofakeit.Sentence(gofakeit.Number(3, 15)),
// 		CreationTimestamp: time.Now().Unix(),
// 		SDuration:         gofakeit.RandomInt([]int{300, 600, 1800, 3600, 86400, 259200}),
// 	}

// 	entityId, err := testData.AddTimeoutForUser(*repoSource, &arrangedTimeout)
// 	if err != nil || entityId == nil {
// 		t.Errorf("Test setup failed: %v", err)
// 	}

// 	cases := []struct {
// 		input          string
// 		expectedOutput *dataModels.Timeout
// 	}{
// 		{idUserWithActiveTimeout, &arrangedTimeout},
// 		{idUserWithNoActiveTimeout, nil},
// 	}

// 	repoToTest := repositories.NewTimeoutsRepository()

// 	// Act & Assert
// 	for _, c := range cases {
// 		// Act
// 		output, err := repoToTest.GetUserTimeout(c.input)
// 		// Assert
// 		if err != nil && c.expectedOutput != nil {
// 			t.Errorf("Timeout expected, error occurred instead: %v", err)
// 		}
// 		if *entityId != c.expectedOutput.Id {
// 			t.Errorf("incorrect output for `Id`: expected `%d` but got `%d`", *entityId, output.Id)
// 		}
// 		if output.CreationTimestamp != c.expectedOutput.CreationTimestamp {
// 			t.Errorf("incorrect output for `CreationTimestamp`: expected `%d` but got `%d`", c.expectedOutput.CreationTimestamp, output.CreationTimestamp)
// 		}
// 		if output.Reason != c.expectedOutput.Reason {
// 			t.Errorf("incorrect output for `Reason`: expected `%s` but got `%s`", c.expectedOutput.Reason, output.Reason)
// 		}
// 		if output.UserId != c.expectedOutput.UserId {
// 			t.Errorf("incorrect output for `UserId`: expected `%s` but got `%s`", c.expectedOutput.UserId, output.UserId)
// 		}
// 		if output.SDuration != c.expectedOutput.SDuration {
// 			t.Errorf("incorrect output for `SDuration`: expected `%d` but got `%d`", c.expectedOutput.SDuration, output.SDuration)
// 		}
// 	}

// 	// Cleanup
// 	testData.RemoveUserTimeout(*repoSource, idUserWithActiveTimeout)

// }
