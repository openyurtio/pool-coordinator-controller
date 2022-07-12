/*
Copyright 2022 The OpenYurt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import "github.com/spf13/pflag"

// CoordinatorControllerOptions is the main settings for the coordinator-controller
type CoordinatorControllerOptions struct {
	MetricsAddr string
	ProbeAddr  string
	EnableLeaderElection    bool

}

// NewCoordinatorControllerOptions creates a new CoordinatorControllerOptions with a default config.
func NewCoordinatorControllerOptions() *CoordinatorControllerOptions {
	return  &CoordinatorControllerOptions{
		MetricsAddr: ":8080",
		ProbeAddr: ":8000",
		EnableLeaderElection: false,
	}
}

// ValidateOptions validates CoordinatorControllerOptions
func (o *CoordinatorControllerOptions)ValidateOptions() error {
	return nil
}

func (o *CoordinatorControllerOptions)AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.MetricsAddr, "metrics-bind-address", o.MetricsAddr, "The address the metric endpoint binds to.")
	fs.StringVar(&o.ProbeAddr, "health-probe-bind-address", o.ProbeAddr, "The address the probe endpoint binds to.")
	fs.BoolVar(&o.EnableLeaderElection, "leader-elect", false, "Enable leader election for controller manager. " +
		"Enabling this will ensure there is only one active controller manager.")
}