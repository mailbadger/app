import React from "react";
import styled from "styled-components";
import { Button, Box } from "grommet";
import { Facebook, Google, Github } from "grommet-icons";

const SocialButton = styled(Button)`
  border: 1px solid ${props => props.color};
  background: ${props => props.color};
  :hover {
    box-shadow: ${props => props.color};
  }
  :focus {
    box-shadow: ${props => props.color};
  }
  color: white;
  ${this} svg {
    fill: white;
    stroke: white;
  }
`;

const SocialButtons = () => (
  <Box direction="column">
    <SocialButton
      color="#4267B2"
      type="button"
      href="/api/auth/facebook"
      icon={<Facebook />}
      primary
      label=" Continue with facebook"
    />
    <SocialButton
      margin={{ top: "small" }}
      color="#4285F4"
      type="button"
      href="/api/auth/google"
      icon={<Google />}
      primary
      label=" Continue with google"
    />
    <SocialButton
      margin={{ top: "small" }}
      color="#333333"
      type="button"
      href="/api/auth/github"
      icon={<Github />}
      primary
      label=" Continue with github"
    />
  </Box>
);

export default SocialButtons;
