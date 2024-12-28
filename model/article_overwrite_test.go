package model

import "testing"





func TestOverwrite(t *testing.T) {
	t.Run("Scenario 1: Update only title of the Article", func(t *testing.T) {
		article := Article{Title: "Old Title", Description: "Initial Description", Body: "Initial Body"}
		newTitle := "New Title"
		expected := Article{Title: newTitle, Description: "Initial Description", Body: "Initial Body"}

		article.Overwrite(newTitle, "", "")

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated the title only while preserving other fields.")
		}
	})

	t.Run("Scenario 2: Update only description of the Article", func(t *testing.T) {
		article := Article{Title: "Initial Title", Description: "Old Description", Body: "Initial Body"}
		newDescription := "New Description"
		expected := Article{Title: "Initial Title", Description: newDescription, Body: "Initial Body"}

		article.Overwrite("", newDescription, "")

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated the description only while preserving other fields.")
		}
	})

	t.Run("Scenario 3: Update only body of the Article", func(t *testing.T) {
		article := Article{Title: "Initial Title", Description: "Initial Description", Body: "Old Body"}
		newBody := "New Body"
		expected := Article{Title: "Initial Title", Description: "Initial Description", Body: newBody}

		article.Overwrite("", "", newBody)

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated the body only while preserving other fields.")
		}
	})

	t.Run("Scenario 4: Update all fields of the Article", func(t *testing.T) {
		article := Article{Title: "Old Title", Description: "Old Description", Body: "Old Body"}
		newTitle := "New Title"
		newDescription := "New Description"
		newBody := "New Body"
		expected := Article{Title: newTitle, Description: newDescription, Body: newBody}

		article.Overwrite(newTitle, newDescription, newBody)

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated all fields of the article.")
		}
	})

	t.Run("Scenario 5: No update when all parameters are empty", func(t *testing.T) {
		article := Article{Title: "Initial Title", Description: "Initial Description", Body: "Initial Body"}
		expected := Article{Title: "Initial Title", Description: "Initial Description", Body: "Initial Body"}

		article.Overwrite("", "", "")

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("No changes to the article when all update parameters are empty.")
		}
	})

	t.Run("Scenario 6: Partial update of non-empty fields only", func(t *testing.T) {
		article := Article{Title: "Old Title", Description: "Old Description", Body: "Initial Body"}
		newTitle := "New Title"
		expected := Article{Title: newTitle, Description: "Old Description", Body: "Initial Body"}

		article.Overwrite(newTitle, "", "")

		if article != expected {
			t.Errorf("Expected %v, got %v", expected, article)
		} else {
			t.Log("Successfully updated non-empty fields only, preserving remaining fields.")
		}
	})
}
