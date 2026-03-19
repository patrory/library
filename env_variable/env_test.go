package env

import (
	"os"
	"testing"
)

func Test_InitEnv(t *testing.T) {
	type args struct {
		serviceName     string
		mandatoryParams []string
		optionalParams  map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		setup   func() error
		cleanup func()
		wantErr bool
	}{
		{
			name: "T1_Success",
			args: args{
				serviceName:     "auth",
				mandatoryParams: []string{"DB_USERNAME", "DB_PASSWORD"},
				optionalParams: map[string]interface{}{
					"refreshAfter": 19,
				},
			},
			setup: func() error {
				// viper uses the service name as an env prefix when AutomaticEnv is enabled.
				// For "auth", the env vars will be looked up as AUTH_DB_USERNAME and AUTH_DB_PASSWORD.
				os.Setenv("AUTH_DB_USERNAME", "user")
				os.Setenv("AUTH_DB_PASSWORD", "pswd")
				return nil
			},
			cleanup: func() {
				os.Unsetenv("AUTH_DB_USERNAME")
				os.Unsetenv("AUTH_DB_PASSWORD")
			},
			wantErr: false,
		},
		{
			name: "T2_Failure",
			args: args{
				serviceName:     "authz",
				mandatoryParams: []string{"DB_URL"},
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			err := InitEnv(tc.args.serviceName, tc.args.mandatoryParams, tc.args.optionalParams)
			if (err != nil) && !tc.wantErr {
				t.Errorf("unexpected error: %v, wantErr: %v", err, tc.wantErr)
			}

			if tc.cleanup != nil {
				t.Cleanup(tc.cleanup)
			}
		})
	}
}
