package main

import (
	"context"
	"fmt"
	"log"
	"os"

	v1 "k8s.io/api/batch/v1"
	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to load in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	gitsync_hash := os.Getenv("GITSYNC_HASH")

	jobName := fmt.Sprintf("quarto-render-%s", gitsync_hash)

	job_namespace := os.Getenv("JOB_NAMESPACE")
	namespace := job_namespace

	job_container_image := os.Getenv("JOB_CONTAINER_IMAGE")

	job_command := os.Getenv("JOB_COMMAND")
	command := []string{"sh", "-c"}
	command = append(command, job_command)

	ttlSecondsAfterFinished := int32(900)
	runAsZero := int64(0)

	job := &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: v1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1core.PodTemplateSpec{
				Spec: v1core.PodSpec{
					SecurityContext: &v1core.PodSecurityContext{
						RunAsUser: &runAsZero,
						FSGroup:   &runAsZero,
					},
					Containers: []v1core.Container{
						{
							Name:    "quarto-render",
							Image:   job_container_image,
							Command: command,
							VolumeMounts: []v1core.VolumeMount{
								{
									Name:      "quarto",
									MountPath: "/quarto",
								},
							},
						},
					},
					Volumes: []v1core.Volume{
						{
							Name: "quarto",
							VolumeSource: v1core.VolumeSource{
								PersistentVolumeClaim: &v1core.PersistentVolumeClaimVolumeSource{
									ClaimName: "blog-pvc",
								},
							},
						},
					},
					RestartPolicy: v1core.RestartPolicyNever,
				},
			},
		},
	}

	_, err = clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Failed to create job: %v", err)
	}

	log.Printf("Job %s created successfully in namespace %s", jobName, namespace)
}
