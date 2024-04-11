package main

type RecipeEntity struct {
	ID           uint             `gorm:"primaryKey"`
	Name         string           `gorm:"not null"`
	Ingredients  string           `gorm:"type:text"` // Liste des ingrédients
	Instructions string           `gorm:"type:text"` // Instructions de préparation
	PrepTime     int              // Temps de préparation en minutes
	CookTime     int              // Temps de cuisson en minutes
	Servings     int              // Nombre de portions
	Feedback     []FeedBackEntity `gorm:"foreignKey:RecipeID"`
	GradeAmount  int              `gorm:"default:0"`
	AverageGrade float64          `gorm:"default:0"`
}

type FeedBackEntity struct {
	ID          uint `gorm:"primaryKey"`
	Grade       uint
	Description uint
	RecipeID    int
}
