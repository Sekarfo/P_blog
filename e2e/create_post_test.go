package main

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestCreatePost(t *testing.T) {
	// Start ChromeDriver
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService("chromedriver", 4444, opts...)
	if err != nil {
		t.Fatalf("Error starting ChromeDriver: %v", err)
	}
	defer service.Stop()

	// Create WebDriver
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		t.Fatalf("Error creating WebDriver: %v", err)
	}
	defer wd.Quit()

	// Navigate to login page
	if err := wd.Get("http://localhost:8080/login.html"); err != nil {
		t.Fatalf("Error navigating to login page: %v", err)
	}

	// Log in
	emailElem, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {
		t.Fatalf("Error finding email input: %v", err)
	}
	passwordElem, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {
		t.Fatalf("Error finding password input: %v", err)
	}
	loginButtonElem, err := wd.FindElement(selenium.ByID, "loginButton")
	if err != nil {
		t.Fatalf("Error finding login button: %v", err)
	}

	emailElem.SendKeys("admin@example.com")
	passwordElem.SendKeys("password123")
	loginButtonElem.Click()

	time.Sleep(2 * time.Second) // Wait for login

	// Navigate to create post page
	if err := wd.Get("http://localhost:8080/posts.html"); err != nil {
		t.Fatalf("Error navigating to posts page: %v", err)
	}

	// Fill out form and submit
	titleElem, err := wd.FindElement(selenium.ByID, "postTitle")
	if err != nil {
		t.Fatalf("Error finding post title input: %v", err)
	}
	contentElem, err := wd.FindElement(selenium.ByID, "postContent")
	if err != nil {
		t.Fatalf("Error finding post content input: %v", err)
	}
	submitButtonElem, err := wd.FindElement(selenium.ByID, "submitPost")
	if err != nil {
		t.Fatalf("Error finding submit button: %v", err)
	}

	titleElem.SendKeys("Test Post")
	contentElem.SendKeys("This is a test post.")
	submitButtonElem.Click()

	time.Sleep(2 * time.Second) // Wait for post creation

	// Verify post in list
	posts, err := wd.FindElements(selenium.ByClassName, "post-title")
	if err != nil {
		t.Fatalf("Error finding posts: %v", err)
	}

	found := false
	for _, post := range posts {
		text, _ := post.Text()
		if text == "Test Post" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Created post not found in the list")
	}
}
