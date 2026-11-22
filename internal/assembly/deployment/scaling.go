package deployment

func Scale(deployment *Deployment, replicas int) {
	if replicas >= 0 {
		deployment.Replicas = replicas
	}
}
