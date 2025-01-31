package main

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestCreatePost(t *testing.T) {
	// ChromeDriver
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService("chromedriver", 4444, opts...)
	if err != nil {
		t.Fatalf("Error starting ChromeDriver: %v", err)
	}
	defer service.Stop()

	// WebDriver
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		t.Fatalf("Error creating WebDriver: %v", err)
	}
	defer wd.Quit()

	// login page
	if err := wd.Get("http://localhost:8080/login.html"); err != nil {
		t.Fatalf("Error navigating to login page: %v", err)
	}

	// login
	emailElem, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {
		t.Fatalf("Error finding email input: %v", err)
	}
	passwordElem, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {
		t.Fatalf("Error finding password input: %v", err)
	}
	loginButtonElem, err := wd.FindElement(selenium.ByID, "submitLogin")
	if err != nil {
		t.Fatalf("Error finding login button: %v", err)
	}

	emailElem.SendKeys("ansarsh1243@gmail.com")
	passwordElem.SendKeys("admin")
	loginButtonElem.Click()

	// wait for login
	wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		_, err := wd.FindElement(selenium.ByID, "homeBtn")
		return err == nil, nil
	}, 5*time.Second)

	// create post page
	if err := wd.Get("http://localhost:8080/posts.html"); err != nil {
		t.Fatalf("Error navigating to posts page: %v", err)
	}

	// find ids and submit keys
	titleElem, err := wd.FindElement(selenium.ByID, "postTitle")
	if err != nil {
		t.Fatalf("Error finding post title input: %v", err)
	}
	contentElem, err := wd.FindElement(selenium.ByID, "postContent")
	if err != nil {
		t.Fatalf("Error finding post content input: %v", err)
	}
	submitButtonElem, err := wd.FindElement(selenium.ByCSSSelector, "button[type='submit']")
	if err != nil {
		t.Fatalf("Error finding submit button: %v", err)
	}

	title := "EXAMPLE TEST TITLE"

	titleElem.SendKeys(title)
	contentElem.SendKeys("This is a test post.")
	submitButtonElem.Click()

	// wait for the new post to appear
	err = wd.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		// check the DOM for the post title
		script := `
			const posts = document.querySelectorAll('#postList li h3 a');
			for (const post of posts) {
				if (post.textContent === arguments[0]) {
					return true;
				}
			}
			return false;
		`
		result, err := wd.ExecuteScript(script, []interface{}{title})
		if err != nil {
			return false, nil // script failed
		}
		return result.(bool), nil
	}, 10*time.Second)

	if err != nil {
		t.Fatalf("Post titled '%s' was not found in the list: %v", title, err)
	}
}
