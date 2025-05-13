package fileserver

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func getMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "File Shaar",
	})
}

func getUploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"title": "File Shaar",
	})
}

func getFiles(c *gin.Context) {
	// read the contents of the files dir
	files, err := os.ReadDir("files")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read files directory",
			"error":   err.Error(),
		})
		return
	}

	// Create a slice to store filenames
	var fileNames []string

	// collect all filenames
	for _, file := range files {
		if !file.IsDir() { // only include files, not dirs
			fileNames = append(fileNames, file.Name())
		}
	}

	// return the filenames as JSON
	c.JSON(http.StatusOK, gin.H{
		"files": fileNames,
	})

}

func getFileByName(c *gin.Context) {
	fileName := c.Param("filename")

	// construct file path
	filePath := filepath.Join("files", fileName)

	// check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "File not found",
		})
		return
	}

	// File exists, serve it
	c.File(filePath)
}

// upload file
func uploadFile(c *gin.Context) {
	// Set a reasonable max file size (e.g., 32 MB)
	c.Request.ParseMultipartForm(32 << 20)

	// Get the file from the form data
	file, header, err := c.Request.FormFile("uploadFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No file uploaded",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	// Ensure the files directory exists
	if err := os.MkdirAll("files", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create upload directory",
			"error":   err.Error(),
		})
		return
	}

	// Create the destination file
	dst := filepath.Join("files", header.Filename)
	out, err := os.Create(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create file",
			"error":   err.Error(),
		})
		return
	}
	defer out.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": header.Filename,
	})
}

func API(LanIp string) {

	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("web/*.html")

	// serve static files
	router.Static("/scripts", "web/scripts")
	router.Static("/styles", "web/styles")

	router.GET("/", getMainPage)
	router.GET("/upload", getUploadPage)

	// get files list
	router.GET("/files", getFiles)

	router.GET("/files/:filename", getFileByName)

	// upload file
	router.POST("/files", uploadFile)

	router.Run(LanIp + ":8080")

}
