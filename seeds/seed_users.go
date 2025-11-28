package main

import (
	"log"

	"github.com/amirkhgraphic/go-arcaptcha-service/initializers"
	"github.com/amirkhgraphic/go-arcaptcha-service/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	sampleUsers := []models.User{
		{Username: "alice", Email: "alice@example.com", Bio: "Product manager and demo user.", Gender: "female", Nationality: "US"},
		{Username: "bob", Email: "bob@example.com", Bio: "Security engineer who loves Go.", Gender: "male", Nationality: "UK"},
		{Username: "carol", Email: "carol@example.com", Bio: "QA specialist focusing on APIs.", Gender: "female", Nationality: "CA"},
		{Username: "dave", Email: "dave@example.com", Bio: "Data analyst testing pagination.", Gender: "male", Nationality: "IR"},
		{Username: "erin", Email: "erin@example.com", Bio: "Designer checking search/filter UX.", Gender: "female", Nationality: "IR"},
		{Username: "frank", Email: "frank@example.com", Bio: "Backend engineer exploring Gorm.", Gender: "male", Nationality: "DE"},
		{Username: "grace", Email: "grace@example.com", Bio: "Mobile dev testing APIs.", Gender: "female", Nationality: "FR"},
		{Username: "heidi", Email: "heidi@example.com", Bio: "Ops specialist running k8s.", Gender: "female", Nationality: "NL"},
		{Username: "ivan", Email: "ivan@example.com", Bio: "Infra engineer automating infra.", Gender: "male", Nationality: "RU"},
		{Username: "judy", Email: "judy@example.com", Bio: "Support lead validating flows.", Gender: "female", Nationality: "BR"},
		{Username: "ken", Email: "ken@example.com", Bio: "Performance tester.", Gender: "male", Nationality: "JP"},
		{Username: "lara", Email: "lara@example.com", Bio: "Fullstack dev in training.", Gender: "female", Nationality: "ES"},
		{Username: "mike", Email: "mike@example.com", Bio: "Database admin.", Gender: "male", Nationality: "US"},
		{Username: "nina", Email: "nina@example.com", Bio: "Frontend engineer.", Gender: "female", Nationality: "SE"},
		{Username: "oscar", Email: "oscar@example.com", Bio: "Cloud architect.", Gender: "male", Nationality: "MX"},
		{Username: "peggy", Email: "peggy@example.com", Bio: "Data scientist.", Gender: "female", Nationality: "AU"},
		{Username: "quentin", Email: "quentin@example.com", Bio: "QA automation.", Gender: "male", Nationality: "CA"},
		{Username: "ruth", Email: "ruth@example.com", Bio: "Product designer.", Gender: "female", Nationality: "UK"},
		{Username: "sam", Email: "sam@example.com", Bio: "API integrator.", Gender: "male", Nationality: "ZA"},
		{Username: "tina", Email: "tina@example.com", Bio: "Security analyst.", Gender: "female", Nationality: "IN"},
	}

	for _, u := range sampleUsers {
		if err := initializers.DB.Where("email = ?", u.Email).FirstOrCreate(&u).Error; err != nil {
			log.Fatalf("failed seeding user %s: %v", u.Email, err)
		}
	}

	log.Printf("seeded %d users (idempotent)", len(sampleUsers))
}
