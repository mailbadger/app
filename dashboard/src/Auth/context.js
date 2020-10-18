import React, { Component } from "react";
import PropTypes from "prop-types";
import { authInstance as axios } from "../axios";

const defaultState = {
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
};

export const AuthContext = React.createContext(defaultState);

class AuthProvider extends Component {
  constructor(props) {
    super(props);

    this.fetchUser = this.fetchUser.bind(this);
    this.setUser = this.setUser.bind(this);
    this.logout = this.logout.bind(this);
    this.clear = this.clear.bind(this);
    this.state = defaultState;
  }

  componentDidMount() {
    this.fetchUser();
  }

  fetchUser() {
    const callApi = async () => {
      try {
        this.setState({ isLoading: true, error: null });
        const result = await axios.get("/api/users/me");
        this.setState({
          error: null,
          isLoading: false,
          isAuthenticated: true,
          user: result.data,
        });
      } catch (error) {
        this.setState({
          error: error.response.data,
          isLoading: false,
          isAuthenticated: false,
          user: null,
        });
      }
    };
    callApi();
  }

  setUser(user) {
    this.setState({
      user: user,
      isAuthenticated: true,
      isLoading: false,
      error: null,
    });
  }

  clear() {
    this.setState(defaultState);
    localStorage.clear();
  }

  logout() {
    const callApi = async () => {
      try {
        this.setState({ isLoading: true, error: null });
        await axios.post("/api/logout");
        this.setState(defaultState);
      } catch (error) {
        this.setState({
          error: error.response.data,
          isLoading: false,
          isAuthenticated: false,
          user: null,
        });
      }
    };
    callApi();
    localStorage.clear();
  }

  render() {
    return (
      <AuthContext.Provider
        value={{
          ...this.state,
          logout: this.logout,
          fetchUser: this.fetchUser,
          setUser: this.setUser,
          clear: this.clear,
        }}
      >
        {this.props.children}
      </AuthContext.Provider>
    );
  }
}

AuthProvider.propTypes = {
  children: PropTypes.element.isRequired,
};

const AuthConsumer = AuthContext.Consumer;

export { AuthProvider, AuthConsumer };
