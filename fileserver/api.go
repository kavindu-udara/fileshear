package fileserver

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

// Get current working directory
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

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
	currentDir := getCurrentDir()

	// read the contents of the current dir
	files, err := os.ReadDir(currentDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to read directory",
			"error":   err.Error(),
		})
		return
	}

	// Create a slice to store filenames
	var fileList []map[string]interface{}

	// collect all file information
	for _, file := range files {
		if !file.IsDir() { // only include files, not dirs
			info, err := file.Info()
			if err != nil {
				continue
			}
			fileList = append(fileList, map[string]interface{}{
				"name":     file.Name(),
				"size":     info.Size(),
				"modified": info.ModTime(),
			})
		}
	}

	// return the filenames as JSON
	c.JSON(http.StatusOK, gin.H{
		"files": fileList,
	})

}

func getFileByName(c *gin.Context) {
	fileName := c.Param("filename")
	currentDir := getCurrentDir()

	// construct file path
	filePath := filepath.Join(currentDir, fileName)

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
	currentDir := getCurrentDir()

	// Parse multipart form with reasonable size limit
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to parse form",
			"error":   err.Error(),
		})
		return
	}

	// Get all files from form data
	files := form.File["uploadFile[]"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No files uploaded",
		})
		return
	}

	uploadedFiles := []string{}

	// Loop through each file and save it
	for _, fileHeader := range files {
		// Create the destination file
		dst := filepath.Join(currentDir, fileHeader.Filename)

		// save the file
		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to save file: " + fileHeader.Filename,
				"error":   err.Error(),
			})
			return
		}
		uploadedFiles = append(uploadedFiles, fileHeader.Filename)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"files":   uploadedFiles,
	})

}

func API(LanIp string) {
	// Define command line flags
	dev := flag.Bool("dev", false, "Run in development mode")
	flag.Parse()

	router := gin.Default()

	// Set web directory based on environment
	var webDir string
	if *dev {
		webDir = "web" // Development path
		gin.SetMode(gin.DebugMode)
	} else {
		webDir = "/usr/local/share/fileshear/web" // Production path
		gin.SetMode(gin.ReleaseMode)
	}

	// Load HTML templates
	router.LoadHTMLGlob(filepath.Join(webDir, "*.html"))

	// serve static files
	router.Static("/scripts", filepath.Join(webDir, "scripts"))
	router.Static("/styles", filepath.Join(webDir, "styles"))

	router.GET("/", getMainPage)
	router.GET("/upload", getUploadPage)

	// get files list
	router.GET("/files", getFiles)

	router.GET("/files/:filename", getFileByName)

	// upload file
	router.POST("/files", uploadFile)

	router.Run(LanIp + ":8080")

}
