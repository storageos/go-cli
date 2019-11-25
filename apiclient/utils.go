package apiclient

import (
	"code.storageos.net/storageos/c2-cli/node"
	"code.storageos.net/storageos/c2-cli/pkg/id"
	"code.storageos.net/storageos/c2-cli/volume"
)

// deploymentsForNode is a utility function returning the list of all
// deployments located on node within volumes.
func deploymentsForNode(uid id.Node, volumes []*volume.Resource) []*node.Deployment {
	var deployments []*node.Deployment
	for _, v := range volumes {
		if v.Master.Node == uid {
			deployments = append(
				deployments,
				&node.Deployment{
					VolumeID:   v.ID,
					Deployment: v.Master,
				},
			)
		}

		for _, r := range v.Replicas {
			if r.Node == uid {
				deployments = append(
					deployments,
					&node.Deployment{
						VolumeID:   v.ID,
						Deployment: r,
					},
				)
			}
		}
	}
	return deployments
}

// mapNodeDeployments builds a mapping from node ID to hosted deployments
// for the list of volumes.
func mapNodeDeployments(volumes []*volume.Resource) map[id.Node][]*node.Deployment {
	deployMap := make(map[id.Node][]*node.Deployment)

	for _, v := range volumes {
		deployMap[v.Master.Node] = append(
			deployMap[v.Master.Node],
			&node.Deployment{
				VolumeID:   v.ID,
				Deployment: v.Master,
			},
		)

		for _, r := range v.Replicas {
			deployMap[r.Node] = append(
				deployMap[r.Node],
				&node.Deployment{
					VolumeID:   v.ID,
					Deployment: r,
				},
			)
		}
	}

	return deployMap
}
