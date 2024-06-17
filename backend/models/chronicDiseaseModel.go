package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DiseaseMetric struct {
	Metric_name   string            `json:"metric_name" validate:"required"`
	Metric_Values map[string]string `json:"metric_values"`
}

type Disease struct {
	Disease_name string          `json:"disease_name" validate:"required"`
	Metrics      []DiseaseMetric `json:"metrics"`
}

type ChronicDisease struct {
	ID        primitive.ObjectID `bson:"_id"`
	User_name *string            `json:"user_name" validate:"required,min=2,max=100"`
	Diseases  []Disease          `json:"diseases"`
}
