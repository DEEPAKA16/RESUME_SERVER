package resume

import (
	"bitresume/config"
	"bitresume/models"
	"database/sql" // Import the sql package to handle NullString
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProjectsData fetches and combines project data from multiple tables.
func GetProjectsData(c *gin.Context) {
	rollno := c.Param("rollno")
	var allProjects []models.Project

	// Step 1: Fetch all base projects for the given rollno from the main 'projects' table.
	// We get the project's unique ID here to use in subsequent queries.
	rows, err := config.DB.Query("SELECT id, title_idea, summary FROM projects WHERE rollno = ?", rollno)
	if err != nil {
		fmt.Println("Error fetching from projects table:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch base project data"})
		return
	}
	defer rows.Close()

	// Step 2: Loop through each base project found.
	for rows.Next() {
		var projectID int
		var title, summary string

		if err := rows.Scan(&projectID, &title, &summary); err != nil {
			fmt.Println("Error scanning base project row:", err)
			continue // Skip this project if there's an error
		}

		// Step 3: For each project, fetch its related data (GitHub link and Tech Stack).

		// Fetch the GitHub link from the 'project_files' table.
		var githubLink sql.NullString // Use sql.NullString to handle potential NULL values
		err := config.DB.QueryRow("SELECT github_link FROM project_files WHERE project_id = ?", projectID).Scan(&githubLink)
		if err != nil && err != sql.ErrNoRows {
			fmt.Println("Error fetching github link for project_id", projectID, ":", err)
			// Decide if you want to skip or continue with an empty link
		}

		// Fetch the list of tech stack names from the 'project_tech_stack' table.
		var techStack []string
		stackRows, err := config.DB.Query("SELECT tech_name FROM project_tech_stack WHERE project_id = ?", projectID)
		if err != nil {
			fmt.Println("Error fetching tech stack for project_id", projectID, ":", err)
			// Continue with an empty stack if there's an error
		} else {
			for stackRows.Next() {
				var techName string
				if err := stackRows.Scan(&techName); err == nil {
					techStack = append(techStack, techName)
				}
			}
			stackRows.Close() // Important to close the inner rows loop
		}

		// Step 4: Combine all fetched data into the final struct.
		project := models.Project{
			Title:       title,
			Description: summary,
			Github:      githubLink.String, // .String provides the value or "" if NULL
			Stack:       techStack,
		}

		allProjects = append(allProjects, project)
	}

	// Final Step: Send the complete, aggregated list to the frontend.
	c.JSON(http.StatusOK, allProjects)
}