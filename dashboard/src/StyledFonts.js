import { createGlobalStyle } from 'styled-components';
import PoppinsMedium from './fonts/Poppins/Poppins-Medium.ttf';
import PoppinsBold from './fonts/Poppins/Poppins-Bold.ttf';

export default createGlobalStyle`
@font-face {
    font-family: "Poppins Medium";
    src: url("${PoppinsMedium}") format('truetype');
    font-weight: normal;
    font-style: normal;
}

@font-face {
    font-family: "Poppins Bold";
    src: url("${PoppinsBold}") format('truetype');
    font-weight: normal;
    font-style: normal;
}
`;
