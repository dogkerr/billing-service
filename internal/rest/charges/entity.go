package charges

import (
	"time"
)

type Charge struct {
	ID                   string    `gorm:"column:id;PRIMARY_KEY;NOT NULL"`
	ContainerID          string    `gorm:"column:container_id;NOT NULL"`
	UserID               string    `gorm:"column:user_id;NOT NULL"`
	TotalCpuUsage        float32   `gorm:"column:total_cpu_usage;NOT NULL"`
	TotalMemoryUsage     float32   `gorm:"column:total_memory_usage;NOT NULL"`
	TotalNetIngressUsage float32   `gorm:"column:total_net_ingress_usage;NOT NULL"`
	TotalNetEgressUsage  float32   `gorm:"column:total_net_egress_usage;NOT NULL"`
	Timestamp            time.Time `gorm:"column:timestamp;NOT NULL"`
	TotalCost            float32   `gorm:"column:total_cost;NOT NULL"`
}

func (m *Charge) TableName() string {
	return "charges_table"
}
