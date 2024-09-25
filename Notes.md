Choice of Framework notes

Backend: Go
Why go:
Go syntax seemed easy to get into with our background. It also was a language we would like to become familiar with.
We will try to use the standard package and not gorilla. We have chosen this to make our app simple and get a strong understanding of the standard package of Go.
Cute mascot. It also seems to be vaguely like C which is a language most of the group is also learning this semester.

Frontend:
Typescript, CSS. HTML. These are languages we know and it will be great for our needs. We are considering a framework like angular but if we don't have to use one we would prefer not to.

Problems with code base:
Makefile doesn't work for any of the group members despite the program being functional when you circumvent it.
Code is written inefficiently with functions being unnecessarily complex and long.
Env variables are in source code instead of as secrets.
Some of the tests seem to be working incorrectly.
md5 is not as safe as it was in 2009.
Missing import in tests.

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
      <td style="padding-bottom: 17px;"><strong>Web Framework Components</strong></td>
      <td style="padding-bottom: 17px;">Pascal case</td>
      <td style="padding-bottom: 17px;"><code>HeaderComponent</code>, <code>UserList</code></td>
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


---
Requirement for Mandatory II
After you have setup a these code quality tools and gone through the issues, your group should create a brief document that answers the following questions:

Do you agree with the findings?

  We agree that it could potentially see our code as having minor problems, but we think it was too early for us to tell if it was something we wanted to shift out attention to at this tage.
  
Which ones did you fix?

<img width="879" alt="Skærmbillede 2024-09-25 kl  16 33 53" src="https://github.com/user-attachments/assets/c3235994-0a84-4c51-811f-f688dff4088e">

<img width="708" alt="Skærmbillede 2024-09-25 kl  16 34 00" src="https://github.com/user-attachments/assets/6ae46597-0231-4a9e-ba1a-17d065d4e275">


  We fixed none of them due to the fact that we wanted to experience bigger conflicts, in order to put it into perspective if it was indeed worth giving our time into these small issues.

Which ones did you ignore?

  We ignored all of the 6/6 issues. We had 1 (1 dot) minor and 5 major (2 dots) issues.

Why?

  We have already answered this.

While only one can setup the integration with SonarQube, everyone should be able to answer the questions above for the exam.
