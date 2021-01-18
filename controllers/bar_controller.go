/*


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
	"fmt"

	testv1alpha1 "github.com/samze/reconciler-runtime-test/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-labs/reconciler-runtime/reconcilers"
)

func BarReconciler(c reconcilers.Config) *reconcilers.ParentReconciler {
	c.Log = c.Log.WithName("BarReconciler")

	return &reconcilers.ParentReconciler{
		Type: &testv1alpha1.Bar{},
		Reconciler: reconcilers.Sequence{
			SecretChildReconciler(c),
			FooReconciler(c),
		},

		Config: c,
	}
}

func SecretChildReconciler(c reconcilers.Config) reconcilers.SubReconciler {
	c.Log = c.Log.WithName("Secret")
	return &reconcilers.ChildReconciler{
		ChildType:     &corev1.Secret{},
		ChildListType: &corev1.SecretList{},

		DesiredChild: func(ctx context.Context, parent *testv1alpha1.Bar) (*corev1.Secret, error) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      parent.Labels,
					Annotations: make(map[string]string),
					Name:        fmt.Sprintf("%s-secret", parent.Name),
					Namespace:   parent.Namespace,
				},
				Data: map[string][]byte{"secret": []byte("123")},
			}
			return secret, nil
		},
		SemanticEquals: func(r1, r2 *corev1.Secret) bool {
			// if the two resources are semantically equal, then we don't need
			// to update the server
			return equality.Semantic.DeepEqual(r1.Data, r2.Data) &&
				equality.Semantic.DeepEqual(r1.Labels, r2.Labels)
		},
		MergeBeforeUpdate: func(actual, desired *corev1.Secret) {
			// mutate actual resource with desired state
			actual.Labels = desired.Labels
			actual.Data = desired.Data
		},
		ReflectChildStatusOnParent: func(parent *testv1alpha1.Bar, child *corev1.Secret, err error) {
			// child is the value of the freshly created/updated/deleted child
			// resource as returned from the api server

			// If a fixed desired resource name is used instead of a generated
			// name, check if the err is because the resource already exists.
			// The ChildReconciler will not claim ownership of another resource.
			//
			// See https://github.com/projectriff/system/blob/1fcdb7a090565d6750f9284a176eb00a3fe14663/pkg/controllers/core/deployer_reconciler.go#L277-L283

			c.Log.Info("status on parent")
		},
		Sanitize: func(child *corev1.Secret) interface{} {
			// log only the resources spec. If the resource contained sensitive
			// values (like a Secret) we'd remove them here so they don't end
			// up in our logs
			return child.GetName()
		},

		Config:     c,
		IndexField: ".metadata.secretController",
	}
}

func FooReconciler(c reconcilers.Config) reconcilers.SubReconciler {
	c.Log = c.Log.WithName("Foo")

	return &reconcilers.SyncReconciler{
		Sync: func(ctx context.Context, parent *testv1alpha1.Bar) error {
			parent.Status.FooStatus = parent.Spec.Foo
			return nil
		},
		Config: c,
	}
}
