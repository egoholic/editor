Feature: Editor signing in
    As an editor
    I want to sign in
    In order to be be able to work on blog content

    Scenario: editor sign-ins with correct credentials
      Given there is a client company
      And it had the following editors:
      When I submited the following credentials:
      | login   | password      |
      | mylogin | mypassword123 |
      Then I got access token

    Scenario: editor sign-ins with wrong credentials
      Given there is a client company
      And it had the following editors:
      When I submited the following credentials:
      | login | password         |
      | login | wrongpassword123 |
      Then I got "unauthorized" response
