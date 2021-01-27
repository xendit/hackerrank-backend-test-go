# hackerrank-backend-test-go

Hi!
Thanks for applying to Xendit. Here we list some important notes for you to work on the test.

### Requirements
- Go 1.11 (Hackerrank Server currently support only for Go 1.11)
- For this test, you need to write your application using SQLite as the database.

### Important Notes

- Do not edit the e2e code. But instead you need to pass all the test written there.
- Please run the test locally first before submitting
  ```bash
  $ make test  # for unit test
  $ make e2e-test #for end to end test (API testing)
  ```
- We use Echo Labstack for the HTTP Router as the default, but you can switch to whatever routing framework you're comfortable to. But please don't change the port number, since it used by Hackerrank by default.
- You can add custom command in the Makefile for self development helper, but don't edit the existing one, since it will be used by Hackerrank. Miss config on the makefile will causing your work won't be graded by Hackerrank.
