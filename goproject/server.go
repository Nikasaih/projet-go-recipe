package main

import (
	"html/template"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize a new HTML homePage
	homePage := template.Must(template.ParseFiles("public/index.html"))
	oneReceipPage := template.Must(template.ParseFiles("public/oneReceip.html"))
	db, err := connectDB()
	if err != nil {
		log.Fatal("error on setup bdd ", err)
	}

	initTable(db)

	if err != nil {
		log.Fatal("error on init table ", err)
	}

	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c *fiber.Ctx) error {
		var receipes []RecipeEntity
		db.Find(&receipes)

		// Injecter les données dans le template
		err := homePage.Execute(c.Response().BodyWriter(), receipes)
		// Verifier qu'il est pas d'erreur
		if err != nil {
			return err
		}

		c.Type("html")
		return nil
	})
	app.Get("/recipe/:recipeID", func(c *fiber.Ctx) error {
		// Extract recipe ID from URL parameter
		recipeID := c.Params("recipeID")

		// Convert recipeID to uint
		id, err := strconv.Atoi(recipeID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "parsing"})
		}

		// Find recipe by ID
		var recipe RecipeEntity

		if err := db.Preload("Feedback").Where("ID = ?", id).First(&recipe).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "404 not found"})
		}

		// Execute template and render HTML
		err = oneReceipPage.Execute(c.Response().BodyWriter(), recipe)
		if err != nil {
			return err
		}

		// Set response type to HTML
		c.Type("html")
		return nil
	})

	app.Get("/add-recipe", func(c *fiber.Ctx) error {
		return c.SendFile("public/recipeForm.html")
	})
	// Define a route for adding a recipe
	app.Post("/add-recipe", func(c *fiber.Ctx) error {
		// Parse form data into the ReceipEntity instance
		var recipe RecipeEntity
		if err := c.BodyParser(&recipe); err != nil {
			log.Println("Error parsing request body:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
		}

		// Add the new recipe to the database
		if err := db.Create(&recipe).Error; err != nil {
			log.Println("Error adding recipe to database:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add recipe"})
		}

		// Respond with a success message
		return c.JSON(fiber.Map{"message": "Recipe added successfully"})
	})

	app.Get("/add-feedback/:receipId", func(c *fiber.Ctx) error {
		return c.SendFile("public/feedbackForm.html")
	})
	// Define a route for adding a recipe
	app.Post("/submit-feedback/:recipeID", func(c *fiber.Ctx) error {
		// Extract recipe ID from URL parameter
		recipeID := c.Params("recipeID")

		// Parse form data into the FeedbackEntity instance
		var feedback FeedBackEntity
		if err := c.BodyParser(&feedback); err != nil {
			log.Println("Error parsing request body:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
		}

		// Set the recipe ID for the feedback
		result, errf := strconv.Atoi(recipeID)
		feedback.RecipeID = result
		if errf != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add feedback"})
		}
		// Add the new feedback to the database
		if err := db.Create(&feedback).Error; err != nil {
			log.Println("Error adding feedback to database:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add feedback"})
		}

		//region compute average grade
		var recipe RecipeEntity
		if err := db.Preload("Feedback").First(&recipe, uint(result)).Error; err != nil {
			log.Fatal("Failed to fetch recipe:", err)
		}
		// Compute average grade
		totalGrades := uint(0)
		numFeedbacks := len(recipe.Feedback)
		for _, feedback := range recipe.Feedback {
			totalGrades += feedback.Grade
		}
		recipe.AverageGrade = float64(totalGrades) / float64(numFeedbacks)
		recipe.GradeAmount = numFeedbacks
		db.Save(&recipe)
		//endregion
		// Respond with a success message
		return c.JSON(fiber.Map{"message": "Feedback added successfully"})
	})
	//fakeContent(db)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3001"))
}
