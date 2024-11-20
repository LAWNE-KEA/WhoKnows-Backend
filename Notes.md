# Choice of Framework notes

# Backend: Go
Why we chose Go:
Go syntax seemed easy to get into with our background. It also was a language we would like to become familiar with.
We will try to use the standard package and not Gorilla. We have chosen this to make our application simple and get a strong understanding of the standard package of the Go framework.
Cute mascot. It also seems to be vaguely like C which is a language most of the group is also learning this semester.

# Frontend:
Typescript, CSS. HTML. These are languages we know and it will be great for our needs. We are considering a framework like Angular but if we don't have to use one we would prefer not to.

# Problems with code base:
  1. Makefile doesn't work for any of the group members despite the program being functional when you circumvent it.
  2. Code is written inefficiently with functions being unnecessarily complex and long.
  3. Env variables are in source code instead of as secrets.
  4. Some of the tests seem to be working incorrectly.
  5. md5 is not as safe as it was in 2009.
  6. Missing import in tests.

Commit conventions, try to keep the format of:

```
<past_tense_verb> + <action>
```

# Conventions that we will follow 
<table style="border-collapse: collapse; width: 100%;">
  <thead>
    <tr>
      <th style="padding-bottom: 17px;">Concept/Context</th>
      <th style="padding-bottom: 17px;">Convention</th>
      <th style="padding-bottom: 17px;">Example</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td style="padding-bottom: 17px;"><strong>Python Variables and Functions</strong></td>
      <td style="padding-bottom: 17px;">Snake case</td>
      <td style="padding-bottom: 17px;"><code>my_variable</code>, <code>my_function</code></td>
    </tr>
    <tr>
      <td style="padding-bottom: 17px;"><strong>Database Tables/Collections</strong></td>
      <td style="padding-bottom: 17px;">Plural</td>
      <td style="padding-bottom: 17px;"><code>customers</code>, <code>orders</code></td>
    </tr>
    <tr>
      <td style="padding-bottom: 17px;"><strong>Relational Database Naming</strong></td>
      <td style="padding-bottom: 17px;">Snake case</td>
      <td style="padding-bottom: 17px;"><code>customer_id</code>, <code>order_date</code></td>
    </tr>
    <tr>
      <td style="padding-bottom: 17px;"><strong>JavaScript Variables and Functions</strong></td>
      <td style="padding-bottom: 17px;">Camel case</td>
      <td style="padding-bottom: 17px;"><code>myVariable</code>, <code>myFunction</code></td>
    </tr>
    <tr>
      <td style="padding-bottom: 17px;"><strong>JavaScript Classes</strong></td>
      <td style="padding-bottom: 17px;">Pascal case</td>
      <td style="padding-bottom: 17px;"><code>MyClass</code>, <code>UserManager</code></td>
    </tr>
    <tr>
      <td style="padding-bottom: 17px;"><strong>CSS Classes</strong></td>
      <td style="padding-bottom: 17px;">Kebab case</td>
      <td style="padding-bottom: 17px;"><code>.my-class</code>, <code>.user-profile</code></td>
    </tr>
    <tr>
      <td style="padding-bottom: 17px;"><strong>Environment Variables</strong></td>
      <td style="padding-bottom: 17px;">Upper snake case</td>
      <td style="padding-bottom: 17px;"><code>DATABASE_URL</code>, <code>API_KEY</code></td>
    </tr>
    <tr>
      <td style="padding-bottom: 17px;"><strong>Constants in Code</strong></td>
      <td style="padding-bottom: 17px;">Upper snake case</td>
      <td style="padding-bottom: 17px;"><code>MAX_SIZE</code>, <code>DEFAULT_COLOR</code></td>
    </tr>
  </tbody>
</table>

---
# Regarding the rewrite
  There were a number of considerations that led to the desition to rewrite the app, I will try to outline them here:

  Until now the backend consisted of one large main file containing all code for the project. This had become unreasonable to work with especially as a team of collaborators. This has been resolved by seperating concerns as seen in  the new project structure.

  The new structure is largely based off what the structure we're used to from other languages:
    whoKnows
    ├── compose.dev.yml
    ├── compose.prod.yml
    ├── LICENSE
    ├── README.md
    └── src
        ├── api
        │   ├── configs
        │   │   └── config.go
        │   ├── handlers
        │   │   ├── login_handler.go
        │   │   ├── logout_handler.go
        │   │   ├── page_handler.go
        │   │   ├── register_handler.go
        │   │   └── search_handler.go
        │   ├── router.go
        │   └── services
        │       ├── helper_services.go
        │       ├── search_service.go
        │       └── user_service.go
        ├── database
        │   ├── database.go
        │   └── seeds
        │       └── seed.json
        ├── Dockerfile.dev
        ├── Dockerfile.prod
        ├── go.mod
        ├── go.sum
        ├── helperTypes
        │   ├── response_data.go
        │   └── weather_response.go
        ├── main.go
        ├── models
        │   ├── page_data.go
        │   ├── search_log.go
        │   ├── token.go
        │   └── user.go
        ├── pages
        │   ├── about.html
        │   ├── login.html
        │   ├── register.html
        │   ├── root.html
        │   ├── search.html
        │   └── weather.html
        ├── schema.sql
        ├── security
        │   ├── jwt.go
        │   └── security.go
        ├── static
        │   └── style.css
        └── tmp
            ├── main
            └── whoknows.db


  During the restructuring the opportunity to refactor and/or rewrite several featutes was taken, as we had become more familiar with golang. During this the desition was made to add Gorm to the project. Gorm made our lives significantly easier and despite the complexity of learning a new framework it ultimately reduced the overall complexity of the backend.

  Another change was the move from MD5 to bcrypt (finally). MD5 hasn't been safe for a while and so it was finally ditched.

  The database was switched to postgresql. We have not worked with postgresql before but seeing as Anders recommended it we decided to give it a shot. It paid off immediately with noticably faster query times, and a better pattern matching from ILIKE etc.
  


---
# Requirement for Mandatory II
After you have setup a these code quality tools and gone through the issues, your group should create a brief document that answers the following questions:

# Do you agree with the findings?

  We agree that it could potentially see our code as having minor problems, but we think it was too early for us to tell if it was something we wanted to shift out attention to at this tage.
  
# Which ones did you fix?

We fixed none of them due to the fact that we wanted to experience bigger conflicts, in order to put it into perspective if it was indeed worth giving our time into these small issues.

Here are some screenshots with issues that Code Climate detected in our repository.

<img width="879" alt="Skærmbillede 2024-09-25 kl  16 33 53" src="https://github.com/user-attachments/assets/c3235994-0a84-4c51-811f-f688dff4088e">

---

<img width="708" alt="Skærmbillede 2024-09-25 kl  16 34 00" src="https://github.com/user-attachments/assets/6ae46597-0231-4a9e-ba1a-17d065d4e275">


# Which ones did you ignore?

  We ignored all of the 6/6 issues. We had 1 (1 dot) minor and 5 major (2 dots) issues.

# Why?

  We have already answered this.

While only one can setup the integration with SonarQube, everyone should be able to answer the questions above for the exam.

# Linting

Consider whether you want to ensure linting before anyone can push their code or make it part of a CI pipeline or both.

  - Because Go implicit have linting capabilities build into it, we consider it as something that we have already been given by Go.


# New sonarqube:
  We have since readded sonarqube and it identified a few issues with the html and logout_handler which we have resolved. They were minor issues and easy to fix.

# GitHub Badge for GitHub Actions

![<TEXT ON SHIELD>](https://github.com/<LAWNE-KEA>/<WhoKnows-Backend>/actions/workflows/<WORKFLOW_FILENAME.yml>/badge.svg?branch=main)
