import { useContext } from "react";
import { withRouter } from "react-router-dom";

import { AuthContext } from "./AuthContext";

const Logout = withRouter(({ history }) => {
  const { logout } = useContext(AuthContext);
  logout();
  history.push("/");

  return null;
});

export default Logout;
