// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ChargesTable struct {
	ID                   uuid.UUID        `json:"id"`
	ContainerID          string           `json:"container_id"`
	UserID               string           `json:"user_id"`
	TotalCpuUsage        float64          `json:"total_cpu_usage"`
	TotalMemoryUsage     float64          `json:"total_memory_usage"`
	TotalNetIngressUsage float64          `json:"total_net_ingress_usage"`
	TotalNetEgressUsage  float64          `json:"total_net_egress_usage"`
	Timestamp            pgtype.Timestamp `json:"timestamp"`
	TotalCost            float64          `json:"total_cost"`
}

type DepositsTable struct {
	ID     uuid.UUID `json:"id"`
	UserID string    `json:"user_id"`
	Value  float64   `json:"value"`
	Status string    `json:"status"`
}

type MutationsTable struct {
	ID        uuid.UUID   `json:"id"`
	UserID    string      `json:"user_id"`
	Mutation  float64     `json:"mutation"`
	Balance   float64     `json:"balance"`
	Type      interface{} `json:"type"`
	DepositID pgtype.UUID `json:"deposit_id"`
	ChargeID  pgtype.UUID `json:"charge_id"`
}
