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
	"fmt"
	"strings"
	"time"

	tdenginev1beta1 "github.com/taosdata/TDengine-Operator/operator/api/v1beta1"
	"github.com/taosdata/TDengine-Operator/operator/controllers/util"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var isTest bool
var testIP string

// TDengineReconciler reconciles a TDengine object
type TDengineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tdengine.operator.taosdata.com,resources=tdengines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tdengine.operator.taosdata.com,resources=tdengines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tdengine.operator.taosdata.com,resources=tdengines/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=storageclasses;statefulsets;persistentvolumeclaims;pods;services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TDengine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *TDengineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	app := &tdenginev1beta1.TDengine{}
	err := r.Get(ctx, req.NamespacedName, app)
	if err != nil {
		if errors.IsNotFound(err) {
			// kubectl delete
			err = r.DeleteService(ctx, req)
			if err != nil {
				logger.Error(err, "delete service error")
			}
			err = r.DeleteStatefulSet(ctx, req)
			if err != nil {
				logger.Error(err, "delete Service error")
			}
			return ctrl.Result{}, nil
		} else {
			return ctrl.Result{}, err
		}
	}
	// service
	wantService := util.ServiceFromApp(app)
	currentService := &corev1.Service{}
	needsCreateService := false
	err = r.Get(ctx, req.NamespacedName, currentService)
	if err != nil {
		if errors.IsNotFound(err) {
			needsCreateService = true
			logger.Info("create TDengine Service")
			err = r.Create(ctx, wantService)
			if err != nil {
				logger.Error(err, "create TDengine Service failed")
				return ctrl.Result{}, err
			}
		} else {
			logger.Error(err, "unexpect error while getting TDengine Service status")
			return ctrl.Result{}, err
		}
	}
	if !needsCreateService {
		err = r.UpdateService(ctx, wantService, currentService)
		if err != nil {
			logger.Error(err, "update TDengine Service failed")
			return ctrl.Result{}, nil
		}
	}

	//stateful set
	wantStatefulSet := util.StatefulSetFromApp(app)
	currentStatefulSet := &appv1.StatefulSet{}
	needsCreateStatefulSet := false
	err = r.Get(ctx, req.NamespacedName, currentStatefulSet)
	if err != nil {
		if errors.IsNotFound(err) {
			// create StatefulSet
			logger.Info("create TDengine StatefulSet")
			err = r.Create(ctx, wantStatefulSet)
			needsCreateStatefulSet = true
			if err != nil {
				logger.Error(err, "create TDengine StatefulSet failed")
				return ctrl.Result{}, err
			}
		} else {
			logger.Error(err, "unexpect error while getting TDengine StatefulSet status")
			return ctrl.Result{}, err
		}
	}

	if !needsCreateStatefulSet {
		err = r.UpdateStatefulSet(ctx, wantStatefulSet, currentStatefulSet)
		if err != nil {
			logger.Error(err, "update TDengine StatefulSet failed")
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

func (r *TDengineReconciler) UpdateStatefulSet(ctx context.Context, want, current *appv1.StatefulSet) error {
	logger := log.FromContext(ctx)
	delta := *current.Spec.Replicas - *want.Spec.Replicas
	if delta > 0 {
		//drop dnode
		//update
		//delete pvc
		ns := client.InNamespace(current.Namespace)
		var pvcs corev1.PersistentVolumeClaimList
		err := r.List(ctx, &pvcs, ns, client.MatchingLabels(current.Spec.Template.Labels))
		if err != nil {
			return err
		}
		needDeletePvcTemplateName := make(map[string]struct{}, len(current.Spec.VolumeClaimTemplates)*int(delta))
		for _, template := range current.Spec.VolumeClaimTemplates {
			for j := 0; j < int(delta); j++ {
				needDeletePvcTemplateName[fmt.Sprintf("%s-%s-%d", template.Name, current.Name, int(*current.Spec.Replicas)-j-1)] = struct{}{}
			}
		}
		deletePvc := make([]corev1.PersistentVolumeClaim, 0, len(current.Spec.VolumeClaimTemplates)*int(delta))
		for _, pvc := range pvcs.Items {
			_, ok := needDeletePvcTemplateName[pvc.Name]
			if ok {
				deletePvc = append(deletePvc, pvc)
			}
		}
		url := fmt.Sprintf("http://%s-0.%s.%s.svc.cluster.local:6041/rest/sql", current.Name, current.Spec.ServiceName, current.Namespace)
		if isTest {
			url = fmt.Sprintf("http://%s:31399/rest/sql", testIP)
		}
		user := "root"
		password := "taosdata"
		dnodeMap, err := util.GetDnodeMap(url, user, password)
		if err != nil {
			return err
		}
		for i := 0; i < int(delta); i++ {
			nodeName := fmt.Sprintf("%s-%d.%s.%s.svc.cluster.local:6030", current.Name, int(*current.Spec.Replicas)-i-1, current.Name, current.Namespace)
			id, exist := dnodeMap[nodeName]
			if !exist {
				continue
			}
			err := util.ExecSql(url, fmt.Sprintf("drop dnode %d", id), user, password)
			if err != nil {
				if !strings.HasSuffix(err.Error(), "Dnode does not exist") {
					return err
				}
			}
			for {
				time.Sleep(100 * time.Millisecond)
				count, err := util.GetDnodeCount(url, user, password)
				if err != nil {
					return err
				}
				if count == int(*current.Spec.Replicas)-i-1 {
					break
				}
			}
		}
		err = r.Update(ctx, want)
		if err != nil {
			return err
		}
		for _, pvc := range deletePvc {
			logger.Info("Deleting PVC", "namespace", pvc.Namespace, "pvc_name", pvc.Name)
			err = r.Delete(ctx, &pvc)
			if err != nil {
				return err
			}
		}
	} else {
		return r.Update(ctx, want)
	}
	return nil
}

func (r *TDengineReconciler) UpdateService(ctx context.Context, want, current *corev1.Service) error {
	return r.Update(ctx, want)
}

func (r *TDengineReconciler) DeleteStatefulSet(ctx context.Context, req ctrl.Request) error {
	currentStatefulSet := &appv1.StatefulSet{}
	err := r.Get(ctx, req.NamespacedName, currentStatefulSet)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	logger := log.FromContext(ctx)
	logger.Info("delete StatefulSet", "namespace", currentStatefulSet.Namespace, "pvc_name", currentStatefulSet.Name)
	return r.Delete(ctx, currentStatefulSet)
}
func (r *TDengineReconciler) DeleteService(ctx context.Context, req ctrl.Request) error {
	currentService := &corev1.Service{}
	err := r.Get(ctx, req.NamespacedName, currentService)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	logger := log.FromContext(ctx)
	logger.Info("delete Service", "namespace", currentService.Namespace, "pvc_name", currentService.Name)
	return r.Delete(ctx, currentService)
}

// SetupWithManager sets up the controller with the Manager.
func (r *TDengineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tdenginev1beta1.TDengine{}).
		Owns(&appv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
