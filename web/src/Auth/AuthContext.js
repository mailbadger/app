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

    this.setSession = this.setSession.bind(this);
  }

  componentDidUpdate() {
    this.checkAuth();
  }

  componentDidMount() {
    this.checkAuth();
  }

  checkAuth() {
    let user = null;
    const token = JSON.parse(localStorage.getItem("token"));
    const isAuthenticated = this.isAuthenticated(token);

    if (isAuthenticated !== this.state.isAuthenticated) {
      if (isAuthenticated) {
        user = JSON.parse(localStorage.getItem("user"));
      }

      this.setState({ isAuthenticated, user, token });
    }
  }

  setSession(data) {
    this.setState(data);
  }

  isAuthenticated(token) {
    return token && new Date().getTime() < token.expires_in;
  }

  render() {
    return (
      <AuthContext.Provider
        value={{ ...this.state, setSession: this.setSession }}
      >
        {this.props.children}
      </AuthContext.Provider>
    );
  }
}

const AuthConsumer = AuthContext.Consumer;

export { AuthProvider, AuthConsumer };
