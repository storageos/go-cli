package types

type CPHealthStatusWithID struct {
	*CPHealthStatus
	ID string
}

type ClusterHealthCP []CPHealthStatusWithID

func (c *ClusterHealthCP) Add(nodeID string, health *CPHealthStatus) {
	*c = append(*c, CPHealthStatusWithID{
		ID:             nodeID,
		CPHealthStatus: health,
	})
}

type DPHealthStatusWithID struct {
	*DPHealthStatus
	ID string
}

type ClusterHealthDP []DPHealthStatusWithID

func (c *ClusterHealthDP) Add(nodeID string, health *DPHealthStatus) {
	*c = append(*c, DPHealthStatusWithID{
		ID:             nodeID,
		DPHealthStatus: health,
	})
}
