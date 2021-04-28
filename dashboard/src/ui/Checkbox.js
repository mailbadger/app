import React from 'react';
import PropTypes from 'prop-types';
import { CheckBox as GrommetCheckBox, Anchor, Box } from 'grommet';

const Checkbox = ({ checked, name, label, handleChange, optionalText }) => (
	<GrommetCheckBox
		onChange={handleChange}
		name={name}
        checked={checked}
		label={
			<Box style={{ color: '#000', display: 'block', fontSize:'16px', lineHeight:'25px' }}>
				{label}˝
				{optionalText && (
					<Anchor as="span" href="" color="#000">
						{' '}
						{optionalText}
					</Anchor>
				)}
			</Box>
		}
	/>
);

Checkbox.propTypes = {
	checked: PropTypes.bool,
	name: PropTypes.string,
	label: PropTypes.string,
    handleChange: PropTypes.func,
	optionalText: PropTypes.string
};

export default Checkbox;
