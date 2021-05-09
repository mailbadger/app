import { React } from 'react';
import { Grid, Box } from 'grommet';
import image from '../images/image.png';
import logo from '../images/logo.png';
import styled from 'styled-components';

const mockData = {
	logo,
	description:
		'The key to impactful  scalable email campaigns is with the badger. Treat your email campaigns like Royalty not Chimps.',
	image
};

const StyledDescription = styled(Box)`
	color: #541388;
	font-size: 20px;
	font-family: 'Poppins Bold';
	margin-top:20px;
`;

const AuthSidebar = () => {
	const { logo, description, image } = mockData;

	return (
		<Grid
			style={{ padding: '84px 38px', height: '100vh', backgroundColor: '#fadcff' }}
			fill
			rows={[ 'auto', '1/2' ]}
			columns={[ 'auto' ]}
			areas={[
				{ name: 'title', start: [ 0, 0 ], end: [ 0, 0 ] },
				{ name: 'image', start: [ 0, 1 ], end: [ 0, 1 ] }
			]}
		>
			<Box gridArea="title">
				<img src={logo} />
				<StyledDescription>{description}</StyledDescription>
			</Box>
			<Box gridArea="image">
				<img style={{ height: '100%' }} src={image} />
			</Box>
		</Grid>
	);
};

export default AuthSidebar;
