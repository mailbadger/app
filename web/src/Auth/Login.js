import React, { Fragment } from "react";
import { Anchor, FormField, Button, TextInput, Paragraph } from "grommet";
import { ErrorMessage } from "formik";
import { Facebook, Google, Github, Mail } from "grommet-icons";

const LoginForm = props => {
  const { handleSubmit, handleChange, errors } = props;

  return (
    <Fragment>
      {errors && errors.message && <div>{errors.message}</div>}
      <form onSubmit={handleSubmit} style={{ width: "90%" }}>
        <FormField label="Email" htmlFor="email">
          <TextInput
            placeholder="you@email.com"
            name="email"
            onChange={handleChange}
          />
          <ErrorMessage name="email" />
        </FormField>
        <FormField label="Password" htmlFor="password">
          <TextInput
            placeholder="****"
            name="password"
            type="password"
            onChange={handleChange}
          />
          <ErrorMessage name="password" />
        </FormField>

        <Anchor href="#" primary label="Forgot my password?" />

        <Button
          icon={<Mail />}
          style={{
            marginTop: "18px",
            marginBottom: "4px",
            background: "#654FAA"
          }}
          type="submit"
          primary
          label="Login with email"
        />

        <Paragraph
          style={{
            borderTop: "1px solid #CACACA",
            marginTop: "14px",
            paddingTop: "0px"
          }}
          size="small"
          textAlign="center"
          alignSelf="center"
          alignContent="center"
        >
          or
        </Paragraph>

        <Button
          style={{
            marginTop: "0px",
            marginBottom: "4px",
            background: "#4267B2"
          }}
          type="button"
          icon={<Facebook />}
          primary
          label=" Login with facebook"
        />
        <Button
          style={{
            marginTop: "4px",
            marginBottom: "4px",
            background: "#4285F4"
          }}
          type="button"
          icon={<Google />}
          primary
          label=" Login with google"
        />
        <Button
          style={{
            marginTop: "4px",
            marginBottom: "4px",
            background: "#333333"
          }}
          type="button"
          icon={<Github />}
          primary
          label=" Login with github"
        />
      </form>
    </Fragment>
  );
};

export default LoginForm;
