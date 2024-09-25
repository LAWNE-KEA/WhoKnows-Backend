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
