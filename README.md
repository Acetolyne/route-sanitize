# route-sanitize
Golang route middleware for the Gin framework, which is a sanitizer that forces developers to sanitize post and optionally get values

## Description
This is a middleware meant to be used with golang webservers. Once setup it forces the developer to sanitize all POST values, and POST values that are not exclusively sanitized will not be available in the route. There is also an option to force sanitization of GET values but this requires the developers to follow a pattern of accessing the GET values thru a specific variable instead of accessing them thru the normal method. Unlike the POST values that can be forced the GET values cannot be removed before hitting the route as they are part of the path.

Any new form fields added to any route will be removed and not available in the route until you add the field into the middleware file. Steps to add new fields are described below.

## Setup

Add the sanitizer.go file with any name ending in .go into your project file.

Add the middleware to all routes automatically. This is typically done as in this example below
```golang
router.Use(SanitizeMiddleware()) //Middlewares for all routes may be a comma seperated list
```

## Explaination of workflow

Under the settings section we can enable sanitization of GET query values if we want. set ```sanitize_get_values``` to true

After sanitization in the route we can access values under the PageValues struct. Access is no harder than it would be getting the values using the normal method but we access the values in a different way. In your route you would need to add the below

```golang
values, exists := c.Get("pagevalues")
if !exists {
    logger.Err.Println("Could not get pagevalues from form sanitizer") //basic error about not being able to access any values
} else {
    postvalues = values.(PageValues).PostValues
    posterrors = values.(PageValues).PostErrors
    getvalues = values.(PageValues).GetValues
    geterrors = values.(PageValues).GetErrors
}
```

After using the above to get the values from the middleware the following will be available

### postvalues 
This is all the posted values that were sanitized, any values not exclusively sanitized are not available to the route. access them using the key name they were assigned to in the middleware. ex postvalues["id"]

### posterrors
This is any errors encounters when sanitizing the values or errors returned from any validators the values were run thru. These are available incase you need to act on the errors instead of logging them or if you want to pass an error message back to the user. This is a slice of errors that can be looped thru incase multiple errors are returned.

### getvalues
This is all the get query values that were sent after being sanitized. Unlike the posted values the get values will still be available in the route using the normal method to retrieve them but those values would not yet be sanitized or validated with any custom validators defined in the middleware since we can't remove these values since they are part of the route itself. Due to this it is up to the developers to create a standard to only access get values from the PageValues struct to ensure the values are sanitized and validated. These values can be accessed using the key name that is in the url query section. ex getvalues["id"]

If there is an error sanitizing or validating the get parameters the request is aborted and you will not reach the requested url.

### geterrors
Similar to the post errors, this field can't be used at the moment as the request will be aborted if the query parameters cant be sanitized. This field is only here for possible expansion.

## POST values

To sanitize POST values navigate to the section in the middleware file labeled ////////POST VALUES

There id a switch statement in this section where you need to add a case statement for any new field. there is an example case statement for a field called id this would allow a form to submit a field named id and the field would then be available to the route as postvalues["id"]. Each new field need to be exclusively added as a case value. below explains each line.

```golang
case "id": //allow the form field named "id"
    id, err := Sanitize(value, "number") //sanitize the value of the field using the sanitizer defined as number
    if err != nil { //add any errors to postvalues to be passed to the route incase you need them for logic or user messages.
        pagevalues.PostErrors = append(pagevalues.PostErrors, err...)
    } else {
        postvalues["id"] = id //If no errors then pass it to the route as the key "id"
    }
    //delete the values from the request so unsanitized values are not available at the route.
    delete(c.Request.PostForm, k)
    delete(c.Request.Form, k)
```

## GET values

If GET parameter sanitization is enabled in the settings section then GET values will also be sanitized and they work in the same way as post values but are still available at the route by the normal access method. As stated above sanitizing the GET parameters requires developers to adopt using a specific way of accessing the values, its no harder than the normal method but unsanitized values will still be available if the developer uses the standard method to access the values. 

If the value can't be sanitized the request is aborted and the user will not get the the requested endpoint.

### Sanitizer
After adding the form field name to the post section you will also need to add a section in the ////SANITIZER section. You will see a switch statement here as well, this coorelates to the second argument in the Sanitize() function. In our above example we used the value "number" in the switch statement in the Sanitize function we need to add a case for each type of sanitization we will pass into the function to tell it to sanitize the value in a certain way.
There is an example that is commented out to show how these statements should be added. each case should set the variable r to a regex statement of allowed characters. Any character not allowed by the regex will be deleted from the value.

If the original value and the sanitized value do not match an error is returned which is handled where the Sanitize function was called.
The errors can be logged instead of printed to stdout so you know what is being sanitized and if you need to update the regex, this makes it easier to debug.

Never remove the default in the switch statement or any field with an invalid type specified when calling the Sanitize method will pass. This would eliminate the sanitized by default behavoir making new form fields instead pass sanitization by default which is not the intent of this middleware.

### Validators
Under the section labeled ////VALIDATORS you can add additional validation you need to perform on the values. The validators will be run right after sanitization. There is an example email validator that is commented out to show how to use it. You need to specify the type of field. The email validator in the example would only be used on field that run ```Sanitize(s, "email")``` 

You can add a sanitizer for as many fields as needed as long as you specify the type. There can also be multiple validators for a single kind

### Additional notes

to use the email validator in the example you need to import the package "net/mail"

all posted values automatically get surrounding whitespace trimmed

in case of multiple fields with the same name, like when a form has two name fields only the first fields value is made available to the route.

Allot of the fmt.Println functions should instead be logged so you can more easily debug issues.

get values use the kind "get" so all get parameters can be sanitized by the same regex. If you need to have different regex patterns or want to add validators you will need to add a switch statement to the GET VALUES section that is similar to the one in the POST VALUES section.

