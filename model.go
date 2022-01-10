package main

// Plan represents a CyVerse subscription plan.
type Plan struct {
	ID                *string            `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	Name              string             `gorm:"not null;unique" json:"name"`
	Description       string             `gorm:"not null" json:"description"`
	PlanQuotaDefaults []PlanQuotaDefault `json:"quota_defaults"`
}

// ResourceType represents a type of resource that a quota may be applied to.
type ResourceType struct {
	ID   *string `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	Name string  `gorm:"not null;unique" json:"name"`
	Unit string  `gorm:"not null;unique" json:"unit"`
}

// PlanQuotaDefault represents a single default quota value for a plan and resource type.
type PlanQuotaDefault struct {
	ID             *string      `gorm:"type:uuid;default:uuid_generate_v1()" json:"id"`
	PlanID         *string      `gorm:"type:uuid;not null" json:"-"`
	ResourceTypeID *string      `gorm:"type:uuid;not null" json:"-"`
	QuotaValue     float64      `gorm:"not null"`
	ResourceType   ResourceType `json:"resource_type"`
}
