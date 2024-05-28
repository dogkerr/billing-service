package charges

type ChargeInput struct {
	ID                   string  `json:"id" binding:"required"`
	UserID               string  `json:"user_id" binding:"required"`
	ContainerID          string  `json:"container_id" binding:"required"`
	TotalCpuUsage        float32 `json:"total_cpu_usage" binding:"required"`
	TotalMemoryUsage     float32 `json:"total_memory_usage" binding:"required"`
	TotalNetIngressUsage float32 `json:"total_net_ingress_usage" binding:"required"`
	TotalNetEgressUsage  float32 `json:"total_net_egress_usage" binding:"required"`
	TotalCost            float32 `json:"total_cost"`
}
