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
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	tdenginev1beta1 "github.com/taosdata/TDengine-Operator/operator/api/v1beta1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
)

func TestAPIs(t *testing.T) {
	before()
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
	Describe("TDengine controller", testTDengineController)
	after()
}

func before() {
	var _ = BeforeSuite(func() {
		logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

		ctx, cancel = context.WithCancel(context.TODO())

		By("bootstrapping test environment")
		testEnv = &envtest.Environment{
			CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
			ErrorIfCRDPathMissing: true,
		}
		useExistingCluster := true
		testEnv.UseExistingCluster = &useExistingCluster

		var err error
		// cfg is defined in this file globally.
		cfg, err = testEnv.Start()
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())

		err = tdenginev1beta1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())

		//+kubebuilder:scaffold:scheme

		k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())
		Expect(k8sClient).NotTo(BeNil())

		k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
			Scheme: scheme.Scheme,
		})
		Expect(err).ToNot(HaveOccurred())

		err = (&TDengineReconciler{
			Client: k8sManager.GetClient(),
			Scheme: k8sManager.GetScheme(),
		}).SetupWithManager(k8sManager)
		Expect(err).ToNot(HaveOccurred())
		s := &storagev1.StorageClassList{}
		err = k8sClient.List(ctx, s)
		if err != nil {
			panic(err)
		}
		if len(s.Items) == 0 {
			createSC()
		}
		err = k8sClient.List(ctx, s)
		if err != nil {
			panic(err)
		}
		go func() {
			defer GinkgoRecover()
			err = k8sManager.Start(ctx)
			Expect(err).ToNot(HaveOccurred(), "failed to run manager")
		}()

	}, 60)
}

func after() {
	var _ = AfterSuite(func() {
		cancel()
		By("tearing down the test environment")
		err := testEnv.Stop()
		Expect(err).NotTo(HaveOccurred())
	})
}

func createSC() {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "local-path-storage",
		},
	}
	err := k8sClient.Create(ctx, ns)
	if err != nil {
		panic(err)
	}

	account := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "local-path-provisioner-service-account",
			Namespace: "local-path-storage",
		},
	}
	err = k8sClient.Create(ctx, account)
	if err != nil {
		panic(err)
	}

	role := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "local-path-provisioner-role",
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"nodes", "persistentvolumeclaims", "configmaps"},
				Verbs:     []string{"get", "list", "watch"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints", "persistentvolumes", "pods"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"events"},
				Verbs:     []string{"create", "patch"},
			},
			{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}
	err = k8sClient.Create(ctx, role)
	if err != nil {
		panic(err)
	}

	roleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "local-path-provisioner-bind",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "local-path-provisioner-role",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "local-path-provisioner-service-account",
				Namespace: "local-path-storage",
			},
		},
	}
	err = k8sClient.Create(ctx, roleBinding)
	if err != nil {
		panic(err)
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "local-path-config",
			Namespace: "local-path-storage",
		},
		Data: map[string]string{
			"config.json": `
    {
            "nodePathMap":[
            {
                    "node":"DEFAULT_PATH_FOR_NON_LISTED_NODES",
                    "paths":["/opt/local-path-provisioner"]
            }
            ]
    }`,
			"setup": `
    #!/bin/sh
    set -eu
    mkdir -m 0777 -p "$VOL_DIR"`,
			"teardown": `
    #!/bin/sh
    set -eu
    rm -rf "$VOL_DIR"`,
			"helperPod.yaml": `
    apiVersion: v1
    kind: Pod
    metadata:
      name: helper-pod
    spec:
      containers:
      - name: helper-pod
        image: busybox
        imagePullPolicy: IfNotPresent`,
		},
	}
	err = k8sClient.Create(ctx, cm)
	if err != nil {
		panic(err)
	}

	replica := int32(1)
	deployment := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "local-path-provisioner",
			Namespace: "local-path-storage",
		},
		Spec: appv1.DeploymentSpec{
			Replicas: &replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "local-path-provisioner",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "local-path-provisioner",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "local-path-provisioner-service-account",
					Containers: []corev1.Container{
						{
							Name:            "local-path-provisioner",
							Image:           "rancher/local-path-provisioner:master-head",
							ImagePullPolicy: "IfNotPresent",
							Command: []string{
								"local-path-provisioner",
								"--debug",
								"start",
								"--config",
								"/etc/config/config.json",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config-volume",
									MountPath: "/etc/config/",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "local-path-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err = k8sClient.Create(ctx, deployment)
	if err != nil {
		panic(err)
	}

	mode := storagev1.VolumeBindingWaitForFirstConsumer
	policy := corev1.PersistentVolumeReclaimDelete
	sc := &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "local-path",
			Annotations: map[string]string{
				"storageclass.kubernetes.io/is-default-class": "true",
			},
		},
		Provisioner:       "rancher.io/local-path",
		VolumeBindingMode: &mode,
		ReclaimPolicy:     &policy,
	}
	err = k8sClient.Create(ctx, sc)
	if err != nil {
		panic(err)
	}

}
