# hackerrank-backend-test-go

Hi!
Thanks for applying to Xendit. Here we list some important notes for you to work on the test.

### Requirements
- Go 1.11 (Hackerrank Server currently support only for Go 1.11)
- For this test, you need to write your application using SQLite as the database.

### Important Notes

- **Do not edit the e2e code**, instead you have to pass all the tests written there.
- Please run the test locally first before submitting
  ```bash
  $ make test  # for unit test
  $ make e2e-test #for end to end test (API testing)
  ```
- We use Echo Labstack for the HTTP Router as the default router, but you can switch to whatever routing framework you're comfortable with.**Don't change the port number as it's Hackerrank's default port.**
- You may add custom command in the Makefile to help you in developing, but **do not** edit the existing ones, as it will be used by Hackerrank. Missing config in the Makefile will cause your work won't be graded by Hackerrank.
