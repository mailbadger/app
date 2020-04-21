import axios from "axios";
import history from "./history";

const xsrfTokenKey = "xsrf_token";

const authInstance = axios.create();

authInstance.interceptors.request.use(
  (config) => {
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
  (err) => {
    return Promise.reject(err);
  }
);

authInstance.interceptors.response.use((res) => {
  const t = res.headers["x-csrf-token"];
  if (t) {
    localStorage.setItem(xsrfTokenKey, t);
  }
  return res;
});

const mainInstance = axios.create();

mainInstance.interceptors.request.use(
  (config) => {
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
  (err) => {
    return Promise.reject(err);
  }
);

mainInstance.interceptors.response.use(
  (res) => {
    const t = res.headers["x-csrf-token"];
    if (t) {
      localStorage.setItem(xsrfTokenKey, t);
    }
    return res;
  },
  (error) => {
    if (error.response && 401 === error.response.status) {
      localStorage.setItem("force_login", "t");
      history.push("/login");
    }

    return Promise.reject(error);
  }
);

export default mainInstance;

export { authInstance, mainInstance };
