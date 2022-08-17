# Deployment to Heroku

> Note : [Docker Desktop](https://www.docker.com/products/docker-desktop/) must be installed on the system deploying the app to Heroku

1. Create Heroku Account
2. Create a project in docker make note of the project name `{Heroku-app-name}`
3. Set `Config Vars`
    1.  Click on the `app-name` and go to `Settings`
    2.  Click on the `Reveal Config Vars`
    3.  Add the following variables (key/value)
    4.  `x-access-token`  : The authorization token that is passed from the heroku service for authorization
    5.  `x-api-token`    : Okta API Key
    6.  `x-org-url`    : Okta Url (Example : https://xyz.oktapreview.com)
4. Open a terminal from the system where the project is docnloaded
5. `cd` to the folder where `Dockerfile` exists
6. execute 
```bash
git init
```
7. execute
```bash
heroku login
```
8. execute
```bash
heroku git:remote -a {Heroku-app-name}
```
9. execute
```bash
heroku container:login
```
10. execute
```bash
heroku container:push web
```
11. execute
```bash
heroku container:release web
```
12. execute
```bash
heroku logs --tail
```

# Test from Postman

1. Import the `UpdateAppHerokuGolangSDK.postman_collection.json` collection to postman
2. Change the `{{heroku-app-name}}` to the app name created at Heroku
3. At the Postman after importing the collection, go to `Headers` and add `x-access-token` with the same value added at the {heroku -> app -> Settings -> Config Vars}
4. At postman, goto `Body` and change the `name` value to {saml app name} and `attributeValue` to the Group attribute statement Value (assuming only one group attr statement) 