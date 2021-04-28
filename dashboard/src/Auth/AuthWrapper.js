import { React } from 'react';
import { ReactComponent as AuthImageBanner } from '../images/AuthImageBanner.svg';
import Login from './Login';
import Register from './Register';
import PropTypes from 'prop-types';
import { Box } from 'grommet';
import { Route, Switch } from 'react-router-dom';

const AuthWrapper = ({ fetchUser }) => (
	<Box direction="row">
		<AuthImageBanner style={{ height: 'auto', width: 'auto' }} />
		<Switch>
			<Route path="/login" component={() => <Login fetchUser={fetchUser} />} />
			<Route path="/signup" component={() => <Register fetchUser={fetchUser} />} />
			<Route path="/" component={() => <Login fetchUser={fetchUser} />} />
		</Switch>
	</Box>
);

AuthWrapper.propTypes = {
	fetchUser: PropTypes.func.isRequired
};

export default AuthWrapper;
