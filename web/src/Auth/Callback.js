import React, { useEffect } from "react";
import queryString from "query-string";
import { withRouter } from "react-router-dom";
import useApi from "../hooks/useApi";

const Callback = withRouter(props => {
  const [state] = useApi({
    url: "/api/users"
  });

  const params = queryString.parse(props.location.search);

  useEffect(() => {
    if (state.data) {
      localStorage.setItem("user", JSON.stringify(state.data));
      props.history.replace("/dashboard");
    }
  }, [state.data, params]);

  const exp = parseInt(params.exp, 10) * 1000 + new Date().getTime();
  localStorage.setItem(
    "token",
    JSON.stringify({ access_token: params.t, expires_in: exp })
  );

  if (!state.data) {
    return <div>Loading...</div>;
  }

  return null;
});

export default Callback;
