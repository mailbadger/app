import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { FormField, Paragraph, Box } from 'grommet';
import { Formik } from 'formik';
import { string, object } from 'yup';
import { NavLink } from 'react-router-dom';
import { mainInstance as axios } from '../axios';
import qs from 'qs';
import SocialButtons from './SocialButtons';
import { socialAuthEnabled } from '../Auth';
import {
	AuthStyledTextInput,
	AuthStyledTextLabel,
	AuthStyledButton,
	AuthStyledHeader,
	AuthErrorMessage,
	AuthStyledRedirectLink,
	CustomLineBreak,
	AuthMainWrapper
} from '../ui';
import { FormPropTypes } from '../PropTypes';

const Form = ({ handleSubmit, handleChange, isSubmitting, errors }) => (
	<Box flex={true} direction="column">
		<AuthStyledRedirectLink text="Don't have an account?" redirectLink="/signup" redirectLabel="Sign up" />
		<Box flex={true} direction="row" alignSelf="center" justify="center" align="center">
			<AuthMainWrapper width="447px">
				<form onSubmit={handleSubmit}>
					<AuthStyledHeader>Welcome back !</AuthStyledHeader>
					<Paragraph
						margin={{ top: '10px', right: '0' }}
						color="#000"
						style={{ fontSize: '23px', lineHeight: '38px', fontFamily: 'Poppins Medium' }}
					>
						We are so excited to see you again !
					</Paragraph>
					<Paragraph textAlign="center" size="small" color="#D85555">
						{errors && errors.message}
					</Paragraph>
					<FormField htmlFor="email" label={<AuthStyledTextLabel>Email</AuthStyledTextLabel>}>
						<AuthStyledTextInput name="email" onChange={handleChange} />
						<AuthErrorMessage name="email" />
					</FormField>
					<FormField htmlFor="password" label={<AuthStyledTextLabel>Password</AuthStyledTextLabel>}>
						<AuthStyledTextInput name="password" type="password" onChange={handleChange} />
						<AuthErrorMessage name="password" />
					</FormField>
					<NavLink
						style={{ color: '#000', fontSize: '13px', fontFamily: 'Poppins Medium' }}
						to="/forgot-password"
					>
						Forgot your password?
					</NavLink>
					<Box>
						<AuthStyledButton
							margin={{ top: 'medium', bottom: 'medium' }}
							disabled={isSubmitting}
							type="submit"
							primary
							label="Login with email"
						/>
					</Box>
					{socialAuthEnabled() && (
						<Fragment>
							<CustomLineBreak text="or" />
							<SocialButtons />
						</Fragment>
					)}
				</form>
			</AuthMainWrapper>
		</Box>
	</Box>
);

Form.propTypes = FormPropTypes;

const loginValidation = object().shape({
	email: string().required('Please enter your email'),
	password: string().required('Please enter your password')
});

const LoginForm = (props) => {
	const handleSubmit = async (values, { setSubmitting, setErrors }) => {
		const callApi = async () => {
			try {
				await axios.post(
					'/api/authenticate',
					qs.stringify({
						username: values.email,
						password: values.password
					})
				);

				setSubmitting(false);

				props.fetchUser();
			} catch (error) {
				setSubmitting(false);
				setErrors(error.response.data);
			}
		};

		await callApi();
	};

	return (
		<Formik initialValues={{ email: '', password: '' }} onSubmit={handleSubmit} validationSchema={loginValidation}>
			{(props) => <Form {...props} />}
		</Formik>
	);
};

LoginForm.propTypes = {
	fetchUser: PropTypes.func.isRequired
};

export default LoginForm;
