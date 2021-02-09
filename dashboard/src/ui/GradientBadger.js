import React from 'react';

import { Box } from 'grommet';
import Emoji from "./Emoji";

const GradientBadger = () => (
  <Box
    background="linear-gradient(#6FFFB0 0%, #7D4CDB 100%)"
    border={{ color: 'white', size: 'small' }}
    margin={{ bottom: 'medium' }}
    pad="xsmall"
    round="small"
  >
    <Emoji symbol="🦡" label="Badger logo" />
  </Box>
);

export default GradientBadger;