package main

import (
	"fmt"

	"github.com/bxcodec/faker/v3" // Importing the faker library
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() (*gorm.DB, error) {
	dsn := "host=postgres user=myuser password=secret dbname=mydatabase port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initTable(db *gorm.DB) {
	db.AutoMigrate(&RecipeEntity{})
	db.AutoMigrate(&FeedBackEntity{})
}

func fakeContent(db *gorm.DB) {
	// Generate fake recipes
	for i := 0; i < 10; i++ {
		outputRecipt := RecipeEntity{}
		// Generate fake data for the recipe using faker library
		err := faker.FakeData(&outputRecipt)
		recipe := RecipeEntity{}
		recipe.Name = outputRecipt.Name
		recipe.Ingredients = outputRecipt.Ingredients
		recipe.Instructions = outputRecipt.Instructions
		recipe.PrepTime = outputRecipt.PrepTime
		recipe.CookTime = outputRecipt.CookTime
		recipe.Servings = outputRecipt.Servings

		if err != nil {
			// Handle the error
			return
		}
		// Save the fake recipe to the database
		db.Create(&recipe)
		fmt.Println(recipe)
	}
}
