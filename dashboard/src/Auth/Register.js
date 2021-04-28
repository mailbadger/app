import React, { useState, Fragment } from 'react';
import PropTypes from 'prop-types';
import { Paragraph, FormField, Box } from 'grommet';
import { Formik } from 'formik';
import ReCAPTCHA from 'react-google-recaptcha';
import { string, object, ref, addMethod, bool } from 'yup';
import { mainInstance as axios } from '../axios';
import qs from 'qs';

import equalTo from '../utils/equalTo';
import { socialAuthEnabled } from '../Auth';
import SocialButtons from './SocialButtons';
import {
	AuthStyledTextInput,
	AuthStyledTextLabel,
	AuthStyledHeader,
	AuthStyledButton,
	AuthErrorMessage,
	CustomLineBreak,
	Checkbox,
	AuthStyledRedirectLink,
	AuthMainWrapper
} from '../ui';
import { FormPropTypes } from '../PropTypes';

addMethod(string, 'equalTo', equalTo);

const registerValidation = object().shape({
	email: string().email('Please enter a valid email').required('Please enter your email'),
	password: string().required('Please enter your password').min(8),
	password_confirm: string()
		.equalTo(ref('password'), "Passwords don't match")
		.required('Confirm Password is required'),
	terms: bool().oneOf([ true ], 'Please accept Terms of Service')
});

const Form = ({ handleSubmit, handleChange, isSubmitting, setFieldValue, errors }) => {
	const [ checked, setChecked ] = useState(false);

	const handleChangeCheckBtn = (e) => {
		setChecked(!checked);
		handleChange(e);
	};

	return (
		<Box flex={true} direction="column">
			<AuthStyledRedirectLink text="Already a member? " redirectLink="/login" redirectLabel="Sign In" />
			<Box flex={true} direction="row" alignSelf="center" justify="center" align="center">
				<AuthMainWrapper width="503px">
					<form onSubmit={handleSubmit}>
						<AuthStyledHeader>Sign up to Mailbadger </AuthStyledHeader>
						<Paragraph textAlign="center" size="small" color="#D85555">
							{errors && errors.message}
						</Paragraph>
						<FormField
							style={{ marginTop: '10px' }}
							htmlFor="email"
							label={<AuthStyledTextLabel>Email</AuthStyledTextLabel>}
						>
							<AuthStyledTextInput name="email" onChange={handleChange} />
							<AuthErrorMessage name="email" />
						</FormField>
						<FormField
							style={{ marginTop: '10px' }}
							htmlFor="password"
							label={<AuthStyledTextLabel>Password</AuthStyledTextLabel>}
						>
							<AuthStyledTextInput name="password" type="password" onChange={handleChange} />
							<AuthErrorMessage name="password" />
						</FormField>
						<FormField
							style={{ marginTop: '10px' }}
							htmlFor="password_confirm"
							label={<AuthStyledTextLabel>Confirm Password</AuthStyledTextLabel>}
						>
							<AuthStyledTextInput name="password_confirm" type="password" onChange={handleChange} />
							<AuthErrorMessage name="password_confirm" />
						</FormField>
						{process.env.REACT_APP_RECAPTCHA_SITE_KEY && (
							<ReCAPTCHA
								sitekey={process.env.REACT_APP_RECAPTCHA_SITE_KEY}
								onChange={(value) => setFieldValue('token_response', value, true)}
							/>
						)}
						<FormField style={{ marginTop: '20px' }} htmlFor="terms">
							<Checkbox
								name="terms"
								label="By clicking any of the Sign Up buttons, I agree to the"
								optionalText="Terms of Service"
								checked={checked}
								handleChange={handleChangeCheckBtn}
							/>
							<AuthErrorMessage name="terms" />
						</FormField>
						<Box>
							<AuthStyledButton
								margin={{ top: 'medium', bottom: 'medium' }}
								disabled={isSubmitting}
								type="submit"
								primary
								label="Sign Up"
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
};

Form.propTypes = FormPropTypes;

const RegisterForm = (props) => {
	const handleSubmit = async (values, { setSubmitting, setErrors }) => {
		const callApi = async () => {
			try {
				let params = {
					email: values.email,
					password: values.password
				};

				if (process.env.REACT_APP_RECAPTCHA_SITE_KEY !== '') {
					if (!values.token_response) {
						setErrors({ message: 'Invalid re-captcha response.' });
						return;
					}

					params.token_response = values.token_response;
				}

				await axios.post('/api/signup', qs.stringify(params));
				props.fetchUser();
			} catch (error) {
				setErrors(error.response.data);
			}
		};

		await callApi();

		//done submitting, set submitting to false
		setSubmitting(false);
	};

	RegisterForm.propTypes = {
		fetchUser: PropTypes.func.isRequired
	};

	return (
		<Formik
			initialValues={{ email: '', password: '', password_confirm: '', terms: false }}
			onSubmit={handleSubmit}
			validationSchema={registerValidation}
		>
			{(props) => <Form {...props} />}
		</Formik>
	);
};

export default RegisterForm;
