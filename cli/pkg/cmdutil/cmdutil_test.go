package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRokuHost(t *testing.T) {
	tests := []struct {
		name        string
		setupViper  func()
		shouldError bool
		errorMsg    string
	}{
		{
			name: "NoHostConfigured",
			setupViper: func() {
				viper.Reset()
			},
			shouldError: true,
			errorMsg:    "no Roku device configured",
		},
		{
			name: "InvalidIPFormat",
			setupViper: func() {
				viper.Reset()
				viper.Set("roku.host", "not-an-ip")
			},
			shouldError: true,
			errorMsg:    "invalid host IP address",
		},
		{
			name: "ValidIPFormat",
			setupViper: func() {
				viper.Reset()
				viper.Set("roku.host", "192.168.1.100")
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupViper()
			helper := &Helper{}

			ip, err := helper.ValidateRokuHost()

			if tt.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Empty(t, ip)
			} else {
				// Note: This will fail if the IP is not reachable
				// For this test, we'll just check that the validation passed
				// without checking connectivity (which would require mocking)
				if err != nil {
					// If there's an error, it should be connectivity-related
					assert.Contains(t, err.Error(), "unable to connect")
				}
			}
		})
	}
}

func TestValidateRokuHost_EmptyString(t *testing.T) {
	viper.Reset()
	viper.Set("roku.host", "")

	helper := &Helper{}
	ip, err := helper.ValidateRokuHost()

	assert.Error(t, err)
	assert.Empty(t, ip)
	assert.Contains(t, err.Error(), "no Roku device configured")
}

func TestAddGroup(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	child1 := &cobra.Command{Use: "child1"}
	child2 := &cobra.Command{Use: "child2"}

	AddGroup(parent, "TestGroup", child1, child2)

	// Verify group was added
	groups := parent.Groups()
	require.Len(t, groups, 1)
	assert.Equal(t, "TestGroup", groups[0].ID)
	assert.Equal(t, "TestGroup", groups[0].Title)

	// Verify children were added with correct group ID
	assert.Equal(t, "TestGroup", child1.GroupID)
	assert.Equal(t, "TestGroup", child2.GroupID)

	// Verify children are in parent's commands
	commands := parent.Commands()
	assert.Contains(t, commands, child1)
	assert.Contains(t, commands, child2)
}

func TestAddGroup_EmptyChildren(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}

	AddGroup(parent, "EmptyGroup")

	groups := parent.Groups()
	require.Len(t, groups, 1)
	assert.Equal(t, "EmptyGroup", groups[0].ID)
}

func TestNewHelper(t *testing.T) {
	// Reset viper to clean state
	viper.Reset()

	helper, err := NewHelper()

	require.NoError(t, err)
	assert.NotNil(t, helper)
}
