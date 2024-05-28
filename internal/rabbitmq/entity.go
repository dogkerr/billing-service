package rabbitmq

type UserMetricsMessage struct {
	ContainerID         string
	UserID              string
	CpuUsage            float32
	MemoryUsage         float32
	NetworkIngressUsage float32
	NetworkEgressUsage  float32
}
