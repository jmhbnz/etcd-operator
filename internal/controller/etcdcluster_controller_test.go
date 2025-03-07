/*
Copyright 2024.

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

package controller

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1alpha1 "go.etcd.io/etcd-operator/api/v1alpha1"
)

func TestControllerReconcile(t *testing.T) {
	ctx := t.Context()

	const resourceName = "test-etcd-cluster"
	typeNamespacedName := types.NamespacedName{
		Name:      resourceName,
		Namespace: "default",
	}

	etcdcluster := &operatorv1alpha1.EtcdCluster{}

	t.Log("Checking for EtcdCluster resource")
	err := k8sClient.Get(ctx, typeNamespacedName, etcdcluster)
	if err != nil && errors.IsNotFound(err) {
		t.Log("Creating EtcdCluster resource")
		resource := &operatorv1alpha1.EtcdCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: "default",
			},
		}
		if createErr := k8sClient.Create(ctx, resource); createErr != nil {
			t.Fatalf("Failed to create EtcdCluster resource: %v", createErr)
		}
	} else if err != nil {
		t.Fatalf("Failed to get EtcdCluster resource: %v", err)
	}

	// Defer a cleanup function to remove the resource after the test finishes.
	defer func() {
		t.Log("Cleaning up the EtcdCluster resource")
		resource := &operatorv1alpha1.EtcdCluster{}
		if getErr := k8sClient.Get(ctx, typeNamespacedName, resource); getErr != nil {
			t.Errorf("Failed to get EtcdCluster resource before deletion: %v", getErr)
		} else {
			if deleteErr := k8sClient.Delete(ctx, resource); deleteErr != nil {
				t.Errorf("Failed to delete EtcdCluster resource: %v", deleteErr)
			}
		}
	}()

	reconciler := &EtcdClusterReconciler{
		Client: k8sClient,
		Scheme: k8sClient.Scheme(),
	}

	_, err = reconciler.Reconcile(ctx, reconcile.Request{
		NamespacedName: typeNamespacedName,
	})
	if err != nil {
		t.Fatalf("Reconciliation failed: %v", err)
	}

	// TODO: Add more specific checks (e.g., verifying status conditions or created resources).
	// For example:
	// updated := &operatorv1alpha1.EtcdCluster{}
	// if err := k8sClient.Get(ctx, typeNamespacedName, updated); err != nil {
	//     t.Fatalf("Failed to retrieve updated resource: %v", err)
	// }
	// // Validate updated fields or status
}
