package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type PageValues struct {
	PostValues map[string]string
	GetValues  map[string]string
	PostErrors []string
	GetErrors  []string
}

/////SETTINGS//////
var sanitize_get_values = false


// Middleware to trim and sanitize all get and post values on all routes
// removes the values from the request and makes them available in the PageValues struct
func FormMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagevalues PageValues
		postvalues := make(map[string]string)
		getvalues := make(map[string]string)
		////////////POST VALUES
		for key, value := range c.Request.Form {
			value := strings.TrimSpace(value[0]) //trim all surrounding space by default and ignore additional values sent by duplicate named fields
			switch k := key; k {
				//Catch all post values here
				case "id":
					id, err := Sanitize(value, "number")
					if err != nil {
						pagevalues.PostErrors = append(pagevalues.PostErrors, err...)
					} else {
						postvalues["id"] = id
					}
					delete(c.Request.PostForm, k)
					delete(c.Request.Form, k)
			}
		}
		//Delete all the values that we didnt sanitize and log a message so we know it needs to be added
		for k := range c.Request.PostForm {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!FIELD NOT SANITIZED!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
			fmt.Println("URI", c.Request.URL.Path, "KEY", k)
			delete(c.Request.PostForm, k)
		}
		for k := range c.Request.Form {
			delete(c.Request.Form, k)
		}
		pagevalues.PostValues = postvalues

		///////GET VALUES IF ENABLED
		if (sanitize_get_values) {
			for k, v := range c.Request.URL.Query() {
				k, err := Sanitize(k, "get")
				if err != nil {
					fmt.Println("Could not sanitize get query key", k, err)
					pagevalues.GetErrors = append(pagevalues.GetErrors, "Problem with query parameters")
					c.Abort()
				} else {
					val, alert := Sanitize(v[0], "get")
					if alert != nil {
						fmt.Println("Could not sanitize get query value", v[0], alert)
						pagevalues.GetErrors = append(pagevalues.GetErrors, "Problem with query parameters")
						c.Abort()
					} else {
						getvalues[k] = val
					}
				}
			}
			pagevalues.GetValues = getvalues
			c.Set("pagevalues", pagevalues)
		}
		c.Next()
	}
}

//////////SANITIZER
func Sanitize(s string, kind string) (updated string, err []string) {
	var r string
	original := s
	//Check for bad characters
	switch k := kind; k {
	// case "number":
	// 	r = `[^\d]+`
	// case "get":
	// 	r = `[a-zA-Z0-9\%]+`
	default:
		r = ``
		err = append(err, "Type", k, "not found in sanitizer")
	}


	rule, error := regexp.Compile(r)
	if err != nil {
		fmt.Println("invalid characters while sanitizing value:", s, "with kind:", kind, error.Error())
		err = append(err, kind+"has invalid characters")
	}
	s = rule.ReplaceAllLiteralString(s, "")
	if s != original {
		err = append(err, kind+" has invalid characters")
		//logs the first character that was stripped so we can update the regex if needed
		for i, char := range s {
			if string(char) != string(original[i]) {
				fmt.Println("[STRIPPED CHAR]", string(original[i]), "from", kind)
				break
			}
		}
	}

	/////////VALIDATORS
	//check for valid email formats
	// if kind == "email" {
	// 	_, error := mail.ParseAddress(s)
	// 	if error != nil {
	// 		err = append(err, "Invalid email address format")
	// 	}
	// 	tld := strings.LastIndex(s, ".")
	// 	if tld < 0 {
	// 		err = append(err, "Invalid email address format")
	// 	}
	// }
	return s, err
}