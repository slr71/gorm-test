package main

import (
	"encoding/json"
	"log"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConfigSpec describes the available configuration settings.
type ConfigSpec struct {
	DatabaseURI string `envconfig:"database_uri" required:"true"`
}

func addResourceType(db *gorm.DB, resourceType *ResourceType) *ResourceType {
	result := db.Create(resourceType)
	log.Printf("addResourceType Result: %+v", result)
	return resourceType
}

func addPlan(db *gorm.DB, plan *Plan) *Plan {
	result := db.Create(plan)
	log.Printf("addPlan Result: %+v", result)
	return plan
}

func addPlanQuotaDefault(db *gorm.DB, planID, resourceTypeID *string, quotaValue float64) {
	planQuotaDefault := &PlanQuotaDefault{
		PlanID:         planID,
		ResourceTypeID: resourceTypeID,
		QuotaValue:     quotaValue,
	}
	result := db.Create(planQuotaDefault)
	log.Printf("addPlanQuotaDefault Result: %+v", result)
}

func getPlan(db *gorm.DB, name string) (*Plan, error) {
	var plan Plan
	result := db.Preload("PlanQuotaDefaults.ResourceType").Where("name = ?", name).First(&plan)
	log.Printf("getPlan Result: %+v", result)
	return &plan, result.Error
}

func main() {
	// Load the configuration.
	var configSpec ConfigSpec
	err := envconfig.Process("qms", &configSpec)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Connect to the database.
	db, err := gorm.Open(postgres.Open(configSpec.DatabaseURI))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize the database.
	db.AutoMigrate(&Plan{})
	db.AutoMigrate(&ResourceType{})
	db.AutoMigrate(&PlanQuotaDefault{})

	// Add some resource types.
	cpuHours := addResourceType(db, &ResourceType{Name: "cpu", Unit: "CPU Hours"})
	storage := addResourceType(db, &ResourceType{Name: "storage", Unit: "Terabytes"})

	// Add some plans.
	free := addPlan(db, &Plan{Name: "free", Description: "Free tier, available to all CyVerse users."})
	premium := addPlan(db, &Plan{Name: "premium", Description: "Premium plan, best for CyVerse power users."})

	// Add some PlanQuotaDefaults for the free tier.
	addPlanQuotaDefault(db, free.ID, cpuHours.ID, 1000)
	addPlanQuotaDefault(db, free.ID, storage.ID, 100)

	// Add some PlanQuotaDefaults for the premium tier.
	addPlanQuotaDefault(db, premium.ID, cpuHours.ID, 10000)
	addPlanQuotaDefault(db, premium.ID, storage.ID, 1000)

	// Look up the free plan for testing.
	retrievedFree, err := getPlan(db, "free")
	if err != nil {
		log.Fatal(err.Error())
	}
	encoded, err := json.MarshalIndent(retrievedFree, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("%s\n", encoded)
}
