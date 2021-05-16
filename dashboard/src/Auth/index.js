import React from 'react';
import { Switch } from 'react-router-dom';
import { AuthConsumer } from './context';
import AuthWrapper from './AuthWrapper';

export const socialAuthEnabled = () => process.env.REACT_APP_ENABLE_SOCIAL_AUTH === 'true';

const Auth = () => {
	return (
		<AuthConsumer>
			{({ fetchUser }) => (
				<Switch>
					<AuthWrapper fetchUser={fetchUser} />
				</Switch>
			)}
		</AuthConsumer>
	);
};

export default Auth;
