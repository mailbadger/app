import React from 'react';
import NewPassword from './NewPassword';
import { Route, Switch } from 'react-router-dom';
import ForgotPassword from './ForgotPassword';
import { AuthConsumer } from './context';
import AuthWrapper from './AuthWrapper';

export const socialAuthEnabled = () => process.env.REACT_APP_ENABLE_SOCIAL_AUTH === 'true';

const Auth = () => {
	return (
		<AuthConsumer>
			{({ fetchUser }) => (
				<Switch>
					<AuthWrapper fetchUser={fetchUser} />
					<Route path="/forgot-password/:token" component={NewPassword} />
					<Route path="/forgot-password" component={ForgotPassword} />
				</Switch>
			)}
		</AuthConsumer>
	);
};

export default Auth;
