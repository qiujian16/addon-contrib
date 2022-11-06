/*
Copyright 2022.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ocmv1alpha1 "open-cluster-management.io/addon-contrib/example-addon/api/v1alpha1"
)

// HelloSpokeReconciler reconciles a HelloSpoke object
type HelloSpokeReconciler struct {
	client.Client
	HubClient   client.Client
	ClusterName string
}

//+kubebuilder:rbac:groups=example.open-cluster-management.io,resources=hellospokes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.open-cluster-management.io,resources=hellospokes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.open-cluster-management.io,resources=hellospokes/finalizers,verbs=update

// Reconcile copys the HelloSpoke from spoke to hub
func (r *HelloSpokeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling HelloSpoke...")
	defer log.Info("done reconciling HelloSpoke")

	var helloSpoke *ocmv1alpha1.HelloSpoke
	err := r.Client.Get(ctx, req.NamespacedName, helloSpoke)
	if err != nil {
		return ctrl.Result{}, err
	}

	hubHelloSpoke := ocmv1alpha1.HelloSpoke{}
	err = r.HubClient.Get(ctx, types.NamespacedName{Namespace: r.ClusterName, Name: helloSpoke.Name}, &hubHelloSpoke)
	switch {
	case errors.IsNotFound(err):
		hubHelloSpoke.Name = helloSpoke.Name
		hubHelloSpoke.Namespace = r.ClusterName
		hubHelloSpoke.Status = helloSpoke.Status
		if err = r.HubClient.Create(ctx, &hubHelloSpoke); err != nil {
			return ctrl.Result{}, err
		}
	case err != nil:
		return ctrl.Result{}, err
	}

	hubHelloSpoke.Status = helloSpoke.Status
	return ctrl.Result{}, r.Client.Update(ctx, &hubHelloSpoke)
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelloSpokeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ocmv1alpha1.HelloSpoke{}).
		Complete(r)
}
