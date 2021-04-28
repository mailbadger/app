import React from 'react';
import PropTypes from 'prop-types';
import { Paragraph } from 'grommet';
import { ErrorMessage } from 'formik';

const AuthErrorMessage = ({ name }) => (
	<Paragraph color="#D85555" margin="none" style={{ marginTop: '5px' }}>
		<ErrorMessage color="#D85555" name={name} />
	</Paragraph>
);
AuthErrorMessage.propTypes = {
	name: PropTypes.string
};

export default AuthErrorMessage;
