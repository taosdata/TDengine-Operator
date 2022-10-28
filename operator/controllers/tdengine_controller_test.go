package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	tdenginev1beta1 "github.com/taosdata/TDengine-Operator/operator/api/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func testTDengineController() {
	const (
		TDName      = "test-tdengine"
		TDNameSpace = "default"

		timeout  = time.Minute * 5
		interval = time.Millisecond * 250
	)
	Context("Should create 3 pods and 3 pvcs.", func() {
		isTest = true
		By("By creating a new TDengine cluster")
		ctx := context.Background()
		replica := int32(3)
		td := &tdenginev1beta1.TDengine{
			TypeMeta: metav1.TypeMeta{
				Kind:       "TDengine",
				APIVersion: "tdengine.operator.taosdata.com/v1beta1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      TDName,
				Namespace: TDNameSpace,
				Labels: map[string]string{
					"app": TDName,
				},
			},
			Spec: tdenginev1beta1.TDengineSpec{
				Replicas:        &replica,
				Image:           "tdengine/tdengine:latest",
				ImagePullPolicy: "Always",
				Env: []corev1.EnvVar{
					{
						Name:  "TZ",
						Value: "Asia/Shanghai",
					},
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "taosdata",
						},
						Spec: corev1.PersistentVolumeClaimSpec{
							AccessModes: []corev1.PersistentVolumeAccessMode{
								corev1.ReadWriteOnce,
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceStorage: resource.MustParse("5Gi"),
								},
							},
						},
					},
				},
			},
		}
		k8sClient.Delete(ctx, td)
		var pvcs corev1.PersistentVolumeClaimList
		ns := client.InNamespace(TDNameSpace)
		err := k8sClient.List(ctx, &pvcs, ns, client.MatchingLabels(map[string]string{
			"app": TDName,
		}))
		if err != nil {
			for _, pvc := range pvcs.Items {
				k8sClient.Delete(ctx, &pvc)
			}
		}
		testNodePortService := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-td-node-port",
				Namespace: TDNameSpace,
				Labels: map[string]string{
					"app": TDName,
				},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Name:     "tcp-6041",
						Protocol: "TCP",
						Port:     6041,
						NodePort: 31399,
					},
				},
				Selector: map[string]string{
					"app": TDName,
				},
				Type: corev1.ServiceTypeNodePort,
			},
		}
		testNodePortKey := types.NamespacedName{Name: "test-td-node-port", Namespace: TDNameSpace}
		existTestNodePortService := &corev1.Service{}
		err = k8sClient.Get(ctx, testNodePortKey, existTestNodePortService)
		if err != nil {
			if errors.IsNotFound(err) {
				err = k8sClient.Create(ctx, testNodePortService)
				if err != nil {
					panic(err)
				}
				//Expect().Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, testNodePortKey, existTestNodePortService)
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
			}
		}
		Expect(k8sClient.Create(ctx, td)).Should(Succeed())
		key := types.NamespacedName{Name: TDName, Namespace: TDNameSpace}
		created := &tdenginev1beta1.TDengine{}
		Eventually(func() bool {
			err := k8sClient.Get(ctx, key, created)
			if err != nil {
				return false
			}
			return true
		}, timeout, interval).Should(BeTrue())
		sst := &appv1.StatefulSet{}
		Eventually(func() bool {
			err := k8sClient.Get(ctx, key, sst)
			if err != nil {
				return false
			}
			return sst.Status.AvailableReplicas == 3
		}, timeout, interval).Should(BeTrue())

		Expect(k8sClient.List(ctx, &pvcs, ns, client.MatchingLabels(sst.Spec.Template.Labels))).Should(Succeed())
		Expect(len(pvcs.Items)).Should(Equal(3))
		Expect(sst.Status.ReadyReplicas).Should(Equal(int32(3)))
		By("By reduce the number of replicas")
		replica = int32(2)
		td.Spec.Replicas = &replica
		Expect(k8sClient.Update(ctx, td)).Should(Succeed())
		Eventually(func() bool {
			err := k8sClient.Get(ctx, key, sst)
			if err != nil {
				return false
			}
			return sst.Status.AvailableReplicas == 2
		}, timeout, interval).Should(BeTrue())
		Eventually(func() bool {
			err := k8sClient.List(ctx, &pvcs, ns, client.MatchingLabels(sst.Spec.Template.Labels))
			if err != nil {
				return false
			}
			return len(pvcs.Items) == 2
		}, timeout, interval)
		By("By delete")
		Expect(k8sClient.Delete(ctx, td)).Should(Succeed())
		Eventually(func() error {
			err = k8sClient.Get(ctx, key, created)
			return err
		}, timeout, interval).Should(HaveOccurred())
		Eventually(func() bool {
			err := k8sClient.Get(ctx, key, sst)
			if errors.IsNotFound(err) {
				return true
			}
			return false
		}, timeout, interval).Should(BeTrue())
		Expect(k8sClient.List(ctx, &pvcs, ns, client.MatchingLabels(map[string]string{
			"app": TDName,
		}))).Should(Succeed())
		Expect(len(pvcs.Items)).Should(Equal(2))
		for _, pvc := range pvcs.Items {
			k8sClient.Delete(ctx, &pvc)
		}
	})
}
