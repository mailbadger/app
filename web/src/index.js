import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import * as serviceWorker from "./serviceWorker";
import axios from "axios";

axios.interceptors.request.use(
  config => {
    const token = JSON.parse(localStorage.getItem("token"));

    if (token && token.access_token) {
      config.headers.Authorization = `Bearer ${token.access_token}`;
    }

    if (
      config.method === "post" ||
      config.method === "put" ||
      config.method === "delete"
    ) {
      config.headers["Content-Type"] = "application/x-www-form-urlencoded";
    }

    return config;
  },
  err => {
    return Promise.reject(err);
  }
);

ReactDOM.render(<App />, document.getElementById("root"));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
