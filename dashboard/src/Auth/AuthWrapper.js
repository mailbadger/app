import { React, useContext } from 'react';
import Login from './Login';
import Register from './Register';
import PropTypes from 'prop-types';
import { Box, Grid, ResponsiveContext } from 'grommet';
import { Route, Switch } from 'react-router-dom';
import ForgotPassword from './ForgotPassword';
import NewPassword from './NewPassword';
import AuthSidebar from './AuthSidebar';

const AuthWrapper = ({ fetchUser }) => {
	const size = useContext(ResponsiveContext);

	const shouldShowSidebar = size !== 'small' && size !== 'medium';
	const columnSize = shouldShowSidebar ? '500px' : 'auto';
	const isMobile = size === 'small' || size === 'medium';

	return (
		<Grid
			fill
			rows={[ 'auto' ]}
			columns={[ columnSize, 'flex' ]}
			areas={[
				{ name: 'image', start: [ 0, 0 ], end: [ 0, 0 ] },
				{ name: 'main', start: [ 1, 0 ], end: [ 1, 0 ] }
			]}
			style={{ height: '100vh' }}
		>
			{shouldShowSidebar && <AuthSidebar fill gridArea="image" style={{ height: '100%', width: '100%' }} />}
			<Box fill gridArea="main" direction="row" style={{ position: 'relative' }}>
				<Switch>
					<Route path="/login" component={() => <Login isMobile={isMobile} fetchUser={fetchUser} />} />
					<Route path="/signup" component={() => <Register isMobile={isMobile} fetchUser={fetchUser} />} />
					<Route path="/forgot-password/:token" component={NewPassword} />
					<Route path="/forgot-password" component={() => <ForgotPassword isMobile={isMobile} />} />
					<Route path="/" component={() => <Login isMobile={isMobile} fetchUser={fetchUser} />} />
				</Switch>
			</Box>
		</Grid>
	);
};

AuthWrapper.propTypes = {
	fetchUser: PropTypes.func.isRequired
};

export default AuthWrapper;
