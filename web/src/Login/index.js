import React, { Component } from "react";
import { Form, FormField, Button, TextInput } from "grommet";

class Login extends Component {
  render() {
    return (
      <Form>
        <FormField name="email" label="Email" required />
        <FormField name="password" label="Password" required type="password" />
        <Button type="submit" primary label="Submit" />
      </Form>
    );
  }
}

export default Login;
