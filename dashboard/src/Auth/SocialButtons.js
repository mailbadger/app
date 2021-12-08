import React from 'react';
import styled from 'styled-components';
import { Button, Box } from 'grommet';
import { Twitter, Google, Github } from 'grommet-icons';
import { endpoints } from '../network/endpoints';

const SocialButton = styled(Button)`
  width: 140px;
  height: 44px;
  border-radius: 20px;
  text-align:center;
  margin-right:15px;
  margin-top:24px;
  border: 1px solid ${(props) => props.color};
  background: ${(props) => props.color};
  :hover {
    box-shadow: ${(props) => props.color};
  }
  :focus {
    box-shadow: ${(props) => props.color};
  }
  color: white;
  ${this} svg {
    fill: white;
    stroke: white;
  }
  :last-of-type {
    margin-right:0;
  }
  svg {
    margin-top: -2px;
  }
`;

const SocialButtons = () => (
	<Box direction="row">
		<SocialButton color="#4285F4" type="button" href={endpoints.signInWithGoogle} icon={<Google />} primary />
		<SocialButton color="#541388" type="button" href={endpoints.signInWithGithub} icon={<Github />} primary />
		<SocialButton color="#000" type="button" href={endpoints.signInWithTwitter} icon={<Twitter />} primary />
	</Box>
);

export default SocialButtons;
