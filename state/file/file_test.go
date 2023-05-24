package filestate_test

import (
	"fmt"
	"os"
	"testing"

	filestate "github.com/duy-ly/nomios-go/state/file"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	filePath = "./checkpoint.nom"
	gtid     = "665ef2f4-b008-4440-b78c-26ba7ce500e6:123"
)

func init() {
	viper.Set("state.file.path", filePath)
}

func Test_NewFileState(t *testing.T) {
	type testCase struct {
		name        string
		before      func()
		after       func()
		expectError bool
	}

	suite := make([]testCase, 0)

	suite = append(suite, testCase{
		name:        "should_success_path_not_exist",
		before:      func() {},
		after:       func() {},
		expectError: false,
	})

	suite = append(suite, testCase{
		name: "should_error_path_is_dir",
		before: func() {
			os.Mkdir(filePath, 0644)
		},
		after: func() {
			os.Remove(filePath)
		},
		expectError: true,
	})

	suite = append(suite, testCase{
		name: "should_error_path_exist_invalid_format",
		before: func() {
			os.Create(filePath)
		},
		after: func() {
			os.Remove(filePath)
		},
		expectError: true,
	})

	suite = append(suite, testCase{
		name: "should_success_path_exist",
		before: func() {
			os.WriteFile(filePath, []byte("{}"), 0644)
		},
		after: func() {
			os.Remove(filePath)
		},
		expectError: false,
	})

	for _, tc := range suite {
		tc.before()

		_, err := filestate.NewFileState()
		if tc.expectError {
			assert.NotEqual(t, nil, err, "tc %s", tc.name)
		} else {
			assert.Equal(t, nil, err, "tc %s", tc.name)
		}

		tc.after()
	}
}

func Test_SaveLastID(t *testing.T) {
	type testCase struct {
		name          string
		before        func()
		after         func()
		path          string
		expectContent string
	}

	suite := make([]testCase, 0)

	suite = append(suite, testCase{
		name:   "should_equal",
		before: func() {},
		after: func() {
			os.Remove(filePath)
		},
		path:          filePath,
		expectContent: fmt.Sprintf(`{"id":"%s"}`, gtid),
	})

	suite = append(suite, testCase{
		name: "should_equal_unchanged",
		before: func() {
			os.WriteFile(filePath, []byte(fmt.Sprintf(`{"id":"%s"}`, gtid)), 0644)
		},
		after: func() {
			os.Remove(filePath)
		},
		path:          filePath,
		expectContent: fmt.Sprintf(`{"id":"%s"}`, gtid),
	})

	for _, tc := range suite {
		tc.before()

		stt, _ := filestate.NewFileState()

		stt.SaveLastID(gtid)

		// validate content
		data, _ := os.ReadFile(tc.path)
		assert.Equal(t, []byte(tc.expectContent), data, "tc %s", tc.name)

		tc.after()
	}
}

func Test_GetLastID(t *testing.T) {
	type testCase struct {
		name     string
		before   func()
		after    func()
		expectID string
	}

	suite := make([]testCase, 0)

	suite = append(suite, testCase{
		name:     "should_empty_file_not_exist",
		before:   func() {},
		after:    func() {},
		expectID: "",
	})

	suite = append(suite, testCase{
		name: "should_empty_file_empty",
		before: func() {
			os.WriteFile(filePath, []byte("{}"), 0644)
		},
		after: func() {
			os.Remove(filePath)
		},
		expectID: "",
	})

	suite = append(suite, testCase{
		name: "should_success_data_exist",
		before: func() {
			os.WriteFile(filePath, []byte(fmt.Sprintf(`{"id":"%s"}`, gtid)), 0644)
		},
		after: func() {
			os.Remove(filePath)
		},
		expectID: gtid,
	})

	for _, tc := range suite {
		tc.before()

		stt, _ := filestate.NewFileState()

		actual := stt.GetLastID()

		// validate content
		assert.Equal(t, tc.expectID, actual, "tc %s", tc.name)

		tc.after()
	}
}
