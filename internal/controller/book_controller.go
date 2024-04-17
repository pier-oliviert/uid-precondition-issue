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
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	authorv1 "tutorial.kubebuilder.io/project/api/v1"
)

// BookReconciler reconciles a Book object
type BookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=author.tutorial.kubebuilder.io,resources=books,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=author.tutorial.kubebuilder.io,resources=books/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=author.tutorial.kubebuilder.io,resources=books/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Book object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *BookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var book authorv1.Book
	if err := r.Client.Get(ctx, req.NamespacedName, &book); err != nil {
		return ctrl.Result{}, err
	}

	if book.Status.Namespace == nil {

		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "break-",
			},
		}
		if err := r.Client.Create(ctx, namespace); err != nil {
			return ctrl.Result{}, err
		}

		book.Status.Namespace = new(string)
		book.Status.Namespace = &namespace.Name
	}

	if book.Status.Chapter == nil {
		chapter := &authorv1.Chapter{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "chapter-",
				Namespace:    *book.Status.Namespace,
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion: book.APIVersion,
					Kind:       book.Kind,
					Name:       book.Name,
					UID:        book.UID,
				}},
			},
		}

		if err := r.Client.Create(ctx, chapter); err != nil {
			return ctrl.Result{}, err
		}

		book.Status.Chapter = new(string)
		*book.Status.Chapter = chapter.Name

	}
	if err := r.Client.Status().Update(ctx, &book); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&authorv1.Book{}).
		Complete(r)
}
