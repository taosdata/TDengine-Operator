package util

import (
	"fmt"

	"github.com/taosdata/TDengine-Operator/operator/api/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StatefulSetFromApp(app *v1beta1.TDengine) *appv1.StatefulSet {
	template := buildTemplate(app)
	ss := &appv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    app.Labels,
		},
		Spec: appv1.StatefulSetSpec{
			ServiceName: app.Name,
			Replicas:    app.Spec.Replicas,
			UpdateStrategy: appv1.StatefulSetUpdateStrategy{
				Type: appv1.OnDeleteStatefulSetStrategyType,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: app.Labels,
			},
			VolumeClaimTemplates: app.Spec.VolumeClaimTemplates,
			Template:             template,
		},
	}
	return ss
}

func ServiceFromApp(app *v1beta1.TDengine) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels:    app.Labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "tcp-6030",
					Protocol: "TCP",
					Port:     6030,
				},
				{
					Name:     "tcp-6041",
					Protocol: "TCP",
					Port:     6041,
				},
				{
					Name:     "tcp-6042",
					Protocol: "TCP",
					Port:     6042,
				},
				{
					Name:     "tcp-6043",
					Protocol: "TCP",
					Port:     6043,
				},
				{
					Name:     "tcp-6044",
					Protocol: "TCP",
					Port:     6044,
				},
				{
					Name:     "tcp-6046",
					Protocol: "TCP",
					Port:     6046,
				},
				{
					Name:     "tcp-6047",
					Protocol: "TCP",
					Port:     6047,
				},
				{
					Name:     "tcp-6048",
					Protocol: "TCP",
					Port:     6048,
				},
				{
					Name:     "tcp-6049",
					Protocol: "TCP",
					Port:     6049,
				},
				{
					Name:     "udp-6044",
					Protocol: "UDP",
					Port:     6044,
				},
				{
					Name:     "udp-6045",
					Protocol: "UDP",
					Port:     6045,
				},
			},
			Selector: app.Labels,
		},
	}
	return service
}

func buildTemplate(app *v1beta1.TDengine) corev1.PodTemplateSpec {
	mergeEnv(app)
	template := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: app.Labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "tdengine",
					Image:           app.Spec.Image,
					ImagePullPolicy: app.Spec.ImagePullPolicy,
					Ports: []corev1.ContainerPort{
						{
							Name:          "tcp-6030",
							Protocol:      "TCP",
							ContainerPort: 6030,
						},
						{
							Name:          "tcp-6041",
							Protocol:      "TCP",
							ContainerPort: 6041,
						},
						{
							Name:          "tcp-6042",
							Protocol:      "TCP",
							ContainerPort: 6042,
						},
						{
							Name:          "tcp-6043",
							Protocol:      "TCP",
							ContainerPort: 6043,
						},
						{
							Name:          "tcp-6044",
							Protocol:      "TCP",
							ContainerPort: 6044,
						},
						{
							Name:          "tcp-6046",
							Protocol:      "TCP",
							ContainerPort: 6046,
						},
						{
							Name:          "tcp-6047",
							Protocol:      "TCP",
							ContainerPort: 6047,
						},
						{
							Name:          "tcp-6048",
							Protocol:      "TCP",
							ContainerPort: 6048,
						},
						{
							Name:          "tcp-6049",
							Protocol:      "TCP",
							ContainerPort: 6049,
						},
						{
							Name:          "udp-6044",
							Protocol:      "UDP",
							ContainerPort: 6044,
						},
						{
							Name:          "udp-6045",
							Protocol:      "UDP",
							ContainerPort: 6045,
						},
					},
					Env: app.Spec.Env,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "taosdata",
							MountPath: "/var/lib/taos",
						},
					},
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{"taos-check"},
							},
						},
						InitialDelaySeconds: 5,
						TimeoutSeconds:      300,
					},
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{"taos-check"},
							},
						},
						InitialDelaySeconds: 5,
						TimeoutSeconds:      20,
					},
					Resources: app.Spec.PodResources,
				},
			},
		},
	}
	return template
}

func mergeEnv(app *v1beta1.TDengine) {
	tmp := make(map[string]int, len(app.Spec.Env)+3)
	for id, envVar := range app.Spec.Env {
		tmp[envVar.Name] = id
	}
	//pod name
	id, exist := tmp["POD_NAME"]
	podNameEnv := corev1.EnvVar{
		Name: "POD_NAME",
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		}}
	if exist {
		app.Spec.Env[id] = podNameEnv
	} else {
		app.Spec.Env = append(app.Spec.Env, podNameEnv)
	}
	//first EP
	id, exist = tmp["TAOS_FIRST_EP"]
	firstEp := corev1.EnvVar{
		Name:  "TAOS_FIRST_EP",
		Value: fmt.Sprintf("%s-0.%s.%s.svc.cluster.local:6030", app.Name, app.Name, app.Namespace),
	}
	if exist {
		app.Spec.Env[id] = firstEp
	} else {
		app.Spec.Env = append(app.Spec.Env, firstEp)
	}
	// fqdn
	id, exist = tmp["TAOS_FIRST_EP"]
	fqdn := corev1.EnvVar{
		Name:  "TAOS_FQDN",
		Value: fmt.Sprintf("$(POD_NAME).%s.%s.svc.cluster.local", app.Name, app.Namespace),
	}
	if exist {
		app.Spec.Env[id] = fqdn
	} else {
		app.Spec.Env = append(app.Spec.Env, fqdn)
	}
}
