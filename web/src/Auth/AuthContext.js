import React, { Component } from "react";

const defaultState = {
  user: null,
  token: null,
  isAuthenticated: false
};

export const AuthContext = React.createContext(defaultState);

class AuthProvider extends Component {
  state = defaultState;

  constructor(props) {
    super(props);

    this.logout = this.logout.bind(this);
    this.setSession = this.setSession.bind(this);
    this.isAuthenticated = this.isAuthenticated.bind(this);
  }

  checkAuth() {
    let user = null;
    const isAuthenticated = this.isAuthenticated();

    if (isAuthenticated !== this.state.isAuthenticated) {
      if (isAuthenticated) {
        user = JSON.parse(localStorage.getItem("user"));
      }

      const token = JSON.parse(localStorage.getItem("token"));
      this.setState({ isAuthenticated, user, token });
    }
  }

  setSession(data) {
    this.setState(data);
  }

  logout() {
    localStorage.clear();
    this.setState(defaultState);
  }

  isAuthenticated() {
    const token = JSON.parse(localStorage.getItem("token"));
    return token && new Date().getTime() < token.expires_in;
  }

  getUser() {
    if (this.state.user) {
      return this.state.user;
    }

    this.checkAuth();
  }

  getToken() {
    if (this.state.token) {
      return this.state.token;
    }

    this.checkAuth();
  }

  render() {
    return (
      <AuthContext.Provider
        value={{
          ...this.state,
          setSession: this.setSession,
          isAuthenticated: this.isAuthenticated,
          getToken: this.getToken,
          getUser: this.getUser,
          logout: this.logout
        }}
      >
        {this.props.children}
      </AuthContext.Provider>
    );
  }
}

const AuthConsumer = AuthContext.Consumer;

export { AuthProvider, AuthConsumer };
