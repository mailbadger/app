import React from 'react';
import styled from 'styled-components';
import { Paragraph } from 'grommet';
import { NavLink } from 'react-router-dom';
import PropTypes from 'prop-types';

const StyledParagraph = styled(Paragraph)`
    padding-top: 10px;
    font-size: 18px;
    line-height: 30px;
    margin-top: 14px;
    margin-right: 25px;
    text-align:center;
    align-self:flex-end;
    align-content:center;
    color: #000;
`;

const AuthStyledRedirectLink = ({ text, redirectLabel, redirectLink }) => (
	<StyledParagraph>
		{`${text} `}
		<NavLink style={{fontFamily: "Poppins Bold", color: '#000'}} to={redirectLink}>{redirectLabel}</NavLink>
	</StyledParagraph>
);

AuthStyledRedirectLink.propTypes = {
	text: PropTypes.string,
	redirectLabel: PropTypes.string,
	redirectLink: PropTypes.string
};

export default AuthStyledRedirectLink;
