import React, { Fragment } from "react";
import { Paragraph, FormField, Button, TextInput, Anchor } from "grommet";
import { ErrorMessage } from "formik";
import { Facebook, Google, Github, Mail } from "grommet-icons";

const RegisterForm = props => {
  const { handleSubmit, handleChange, errors } = props;

  return (
    <Fragment>
      {errors && errors.message && <div>{errors.message}</div>}
      <form
        onSubmit={handleSubmit}
        style={{ color: "black", width: "90%", height: "100%" }}
      >
        <FormField label="Email" htmlFor="email">
          <TextInput
            style={{ border: "1px solid #CACACA" }}
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
        <FormField label="Confirm Password" htmlFor="confirm-password">
          <TextInput
            placeholder="****"
            name="password"
            type="password"
            onChange={handleChange}
          />
          <ErrorMessage name="password" />
        </FormField>

        <Paragraph
          style={{ marginTop: "14px", paddingTop: "0px" }}
          size="small"
          textAlign="center"
          alignSelf="center"
          alignContent="center"
        >
          By clicking any of the Sign Up buttons, I agree to the{" "}
          <Anchor href="">terms of service</Anchor>
        </Paragraph>

        <Button
          plain
          style={{
            marginTop: "10px",
            marginBottom: "10px",
            borderRadius: "5px",
            padding: "8px",
            background: "#654FAA",
            width: "100%",
            textAlign: "center"
          }}
          icon={<Mail />}
          type="submit"
          alignSelf="stretch"
          textAlign="center"
          primary
          label="Sign UP with email"
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
          plain
          style={{
            marginTop: "7px",
            marginBottom: "7px",
            borderRadius: "5px",
            padding: "8px",
            background: "#4267B2",
            width: "100%",
            textAlign: "center"
          }}
          type="button"
          icon={<Facebook />}
          alignSelf="center"
          textAlign="center"
          primary
          label=" Sign Up with facebook"
        />
        <br />
        <Button
          plain
          style={{
            marginTop: "7px",
            marginBottom: "7px",
            borderRadius: "5px",
            padding: "8px",
            background: "#4285F4",
            width: "100%",
            textAlign: "center"
          }}
          type="button"
          icon={<Google />}
          alignSelf="stretch"
          textAlign="center"
          primary
          label=" Sign Up with google"
        />
        <br />
        <Button
          plain
          style={{
            marginTop: "7px",
            marginBottom: "7px",
            borderRadius: "5px",
            padding: "8px",
            background: "#333333",
            width: "100%",
            textAlign: "center"
          }}
          type="button"
          icon={<Github />}
          alignSelf="stretch"
          textAlign="center"
          primary
          label=" Sign Up with github"
        />
      </form>
    </Fragment>
  );
};

export default RegisterForm;
