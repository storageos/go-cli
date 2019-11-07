package apiclient

import (
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/pkg/volume"
)

// TODO: This is a temp file - should live somewhere else - maybe not be funcs?

func deploymentsForNode(node id.Node, volumes []*volume.Resource) []*volume.Deployment {
	var deployments []*volume.Deployment
	for _, v := range volumes {
		if v.Master.Node == node {
			deployments = append(deployments, v.Master)
		}

		for _, r := range v.Replicas {
			if r.Node == node {
				deployments = append(deployments, v.Master)
			}
		}
	}
	return deployments
}

func mapNodeDeployments(volumes []*volume.Resource) map[id.Node][]*volume.Deployment {
	deployMap := make(map[id.Node][]*volume.Deployment)

	for _, v := range volumes {
		deployMap[v.Master.Node] = append(deployMap[v.Master.Node], v.Master)

		for _, r := range v.Replicas {
			deployMap[r.Node] = append(deployMap[r.Node], r)
		}
	}

	return deployMap
}
