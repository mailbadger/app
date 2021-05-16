import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Router, Route, Switch } from 'react-router-dom';
import { ErrorBoundary } from 'react-error-boundary';
import { Box, Grommet, Heading, Paragraph } from 'grommet';
import { grommet } from 'grommet/themes';
import { deepMerge } from 'grommet/utils';

import { Emoji } from './ui';
import { AuthProvider } from './Auth/context';
import Landing from './Landing';
import Dashboard from './Dashboard';
import Logout from './Auth/Logout';
import VerifyEmail from './VerifyEmail';
import ProtectedRoute from './ProtectedRoute';
import history from './history';

import GlobalFonts from './StyledFonts';

const theme = deepMerge(grommet, {
	global: {
		font: {
			family: "'Poppins', sans-serif;"
		},
		colors: {
			background: '#ffffff',
			brand: '#390099'
		},
		size: {
			avatar: '45px'
		},
		breakpoints: {
			small: { value: 768 },
			medium: {value: 1025}
		}
	},
	tabs: {
		header: {
			background: 'white'
		}
	},
	tab: {
		color: '#888888',
		active: {
			color: 'brand'
		},
		border: false
	},
	formField: {
		label: {
			color: '#ACACAC',
			size: 'small'
		},
		border: false,
		borderColor: '#CACACA',
		margin: 0
	},
	button: {
		border: {
			radius: '5px',
			color: '#390099'
		},
		primary: {
			color: '#390099'
		}
	},
	anchor: {
		primary: {
			color: '#999999'
		},
		color: '#390099'
	},

	checkBox: {
		border: {
			color: '#fadcff'
		},
		size: '30px',
		check: {
			extend: `
			background:#f5f5fa;
			box-shadow: 0 0 1px #390099;
			border-radius: 5px;
			`
		}
	}
});

const ErrorFallback = ({ error }) => (
	<Box align="center" margin={{ top: '10%' }}>
		<Heading level="2">
			Sorry, something went wrong <Emoji label="frowny-face" symbol="😣" />
		</Heading>
		<Box align="start">
			<Paragraph>
				You can open up an issue&nbsp;
				<a
					target="_blank"
					rel="noopener noreferrer"
					href="https://github.com/mailbadger/app/issues/new?assignees=&labels=&template=bug_report.md&title="
				>
					here.
				</a>
				&nbsp;In the meantime, try and reload the page.
			</Paragraph>
			<Paragraph>
				<strong>Error:</strong> {error.message}
			</Paragraph>
		</Box>
	</Box>
);

ErrorFallback.propTypes = {
	error: PropTypes.shape({
		message: PropTypes.string
	})
};

class App extends Component {
	render() {
		return (
			<Grommet theme={theme} full>
				<GlobalFonts />
				<Router history={history}>
					<AuthProvider>
						<Box fill background="background">
							<ErrorBoundary FallbackComponent={ErrorFallback}>
								<Switch>
									<ProtectedRoute path="/dashboard" component={Dashboard} />
									<Route path="/logout" component={Logout} />
									<Route path="/verify-email/:token" component={VerifyEmail} />
									<Route path="/" component={Landing} />
								</Switch>
							</ErrorBoundary>
						</Box>
					</AuthProvider>
				</Router>
			</Grommet>
		);
	}
}

export default App;
