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

package app

import (
	appsv1alpha1 "github.com/openyurtio/coordinator-controller/api/v1beta1"
	"github.com/openyurtio/coordinator-controller/cmd/coordinator-controller/options"
	"github.com/openyurtio/coordinator-controller/controllers"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"os"
	"time"
)

var (
	scheme = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appsv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

// NewCmdCoordinatorController creates a *cobra.Command object with default parameters
func NewCmdCoordinatorController(stopCh <-chan struct{}) *cobra.Command {
	o := options.NewCoordinatorControllerOptions()
	cmd := &cobra.Command{
		Use: "coordinator-controller",
		Short: "Launch coordinator-controller",
		Long: "Launch coordinator-controller",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				klog.V(1).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
			})

			if err := o.ValidateOptions(); err != nil {
				klog.Fatalf("validate options: %v", err)
			}
			Run(o, stopCh)
		},
	}
	o.AddFlags(cmd.Flags())
	return cmd
}

func Run(o *options.CoordinatorControllerOptions, stopCh <-chan struct{})  {
	ctrl.SetLogger(klogr.New())

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     o.MetricsAddr,
		//Port:                   9443,
		HealthProbeBindAddress: o.ProbeAddr,
		LeaderElection:         o.EnableLeaderElection,
		LeaderElectionID:       "coordinator-controller",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	setupLog.Info("[preflight] Running pre-flight checks")
	if err := wait.PollImmediateUntil(time.Second, func() (done bool, err error) {
		// check whether yurtappset CRD exists

		return false, nil
	}, stopCh); err != nil{
		setupLog.Error(err, "failed to wait for discovery yurtappset")
		os.Exit(1)
	}

	setupLog.Info("[prepare] Running prepare for pool coordinator")
	if err := prepareForPoolCoordinator(); err != nil {
		setupLog.Error(err, "failed to prepare for pool coordinator")
		os.Exit(1)
	}

	// setup the coordinator-controller Reconciler and Syncer
	if err = (&controllers.PoolCoordinatorReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "PoolCoordinator")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func prepareForPoolCoordinator() error {
	// 1.prepares the client certificate to access the kubelet for pool-coordinator,
	// saves the certificate in secret and mounts it to pool-coordinator.


	// 2.prepares the client certificate for forwarding node lease to cloud by yurthub,
	// saves the client certificate in secret and will be used by leader yurthub.

	// 3.creates service for pool-coordinator.
	return nil
}