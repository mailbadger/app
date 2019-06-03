import React, { Fragment } from "react";
import { Button } from "grommet";
import { Facebook, Google, Github } from "grommet-icons";

const SocialButtons = () => (
  <Fragment>
    <Button
      style={{
        marginTop: "4px",
        marginBottom: "4px",
        background: "#4267B2"
      }}
      type="button"
      href="/api/auth/facebook"
      icon={<Facebook />}
      primary
      label=" Continue with facebook"
    />
    <br />
    <Button
      style={{
        marginTop: "4px",
        marginBottom: "4px",
        background: "#4285F4"
      }}
      type="button"
      href="/api/auth/google"
      icon={<Google />}
      primary
      label=" Continue with google"
    />
    <br />
    <Button
      style={{
        marginTop: "4px",
        marginBottom: "4px",
        background: "#333333"
      }}
      type="button"
      href="/api/auth/github"
      icon={<Github />}
      primary
      label=" Continue with github"
    />
  </Fragment>
);

export default SocialButtons;
