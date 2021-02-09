import React from 'react';
import PropTypes from 'prop-types';
import { Box } from 'grommet';

const Avatar = ({ name, ...rest }) => (
  <Box
    alignContent="center"
    a11yTitle={`${name} avatar`}
    height="avatar"
    width="avatar"
    round="full"
    background="url(https://www.gravatar.com/avatar/94d093eda664addd6e450d7e9881bcad?s=80&d=identicon)"
    {...rest}
  />
);

Avatar.propTypes = {
  name: PropTypes.string,
};

export default Avatar;