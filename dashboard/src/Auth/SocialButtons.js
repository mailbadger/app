import React from "react";
import { Button, Box } from "grommet";
import { Facebook, Google, Github } from "grommet-icons";

const SocialButtons = () => (
  <Box direction="column">
    <Button
      style={{
        background: "#4267B2"
      }}
      type="button"
      href="/api/auth/facebook"
      icon={<Facebook />}
      primary
      label=" Continue with facebook"
    />
    <Button
      margin={{ top: "small" }}
      style={{
        background: "#4285F4"
      }}
      type="button"
      href="/api/auth/google"
      icon={<Google />}
      primary
      label=" Continue with google"
    />
    <Button
      margin={{ top: "small" }}
      style={{
        background: "#333333"
      }}
      type="button"
      href="/api/auth/github"
      icon={<Github />}
      primary
      label=" Continue with github"
    />
  </Box>
);

export default SocialButtons;
