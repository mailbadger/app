import React from 'react';
import PropTypes from 'prop-types';
import { Menu, Text } from 'grommet';
import Avatar from './Avatar';

const UserMenu = ({ items = [], ...rest }) => (
  <Menu
    dropAlign={{ top: 'bottom', right: 'right' }}
    icon={false}
    items={items.map(item => ({
      ...item,
      label: <Text size="small">{item.label}</Text>,
      onClick: () => {}, // no-op
    }))}
    label={<Avatar />}
    {...rest}
  />
);

UserMenu.propTypes = {
  items: PropTypes.array,
}

export default UserMenu;