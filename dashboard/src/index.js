import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import * as serviceWorker from "./serviceWorker";
import axios from "axios";
import history from "./history";

const xsrfTokenKey = "xsrf_token";

axios.interceptors.request.use(
  config => {
    config.withCredentials = true;

    if (
      config.method === "post" ||
      config.method === "put" ||
      config.method === "patch" ||
      config.method === "delete"
    ) {
      config.headers["Content-Type"] = "application/x-www-form-urlencoded";
      const t = localStorage.getItem(xsrfTokenKey);
      if (t) {
        config.headers["X-CSRF-Token"] = t;
      }
    }

    return config;
  },
  err => {
    return Promise.reject(err);
  }
);

axios.interceptors.response.use(
  res => {
    const t = res.headers["x-csrf-token"];
    if (t) {
      localStorage.setItem(xsrfTokenKey, t);
    }
    return res;
  },
  error => {
    if (error.response && 401 === error.response.status) {
      history.push("/login");
    }

    return Promise.reject(error);
  }
);

ReactDOM.render(<App />, document.getElementById("root"));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
