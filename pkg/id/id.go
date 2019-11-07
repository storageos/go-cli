package id

type Cluster string

func (c Cluster) String() string {
	return string(c)
}

type Node string

func (n Node) String() string {
	return string(n)
}

type Volume string

func (v Volume) String() string {
	return string(v)
}

type Deployment string

func (d Deployment) String() string {
	return string(d)
}

type Namespace string

func (n Namespace) String() string {
	return string(n)
}
