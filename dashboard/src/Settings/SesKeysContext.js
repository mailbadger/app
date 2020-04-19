import React, { Component } from "react";
import PropTypes from "prop-types";
import axios from "axios";

const defaultState = {
  keys: null,
  isLoading: false,
  error: null,
};

export const SesKeysContext = React.createContext(defaultState);

class SesKeysProvider extends Component {
  constructor(props) {
    super(props);

    this.fetchKeys = this.fetchKeys.bind(this);
    this.setKeys = this.setKeys.bind(this);
    this.state = defaultState;
  }

  componentDidMount() {
    this.fetchKeys();
  }

  fetchKeys() {
    const callApi = async () => {
      try {
        this.setState({ isLoading: true, error: null });
        const result = await axios.get("/api/ses-keys");
        this.setState({
          error: null,
          isLoading: false,
          keys: result.data,
        });
      } catch (error) {
        this.setState({
          error: error.response.data,
          isLoading: false,
          keys: null,
        });
      }
    };
    callApi();
  }

  setKeys(keys) {
    this.setState({
      keys: keys,
      isLoading: false,
      error: null,
    });
  }

  render() {
    return (
      <SesKeysContext.Provider
        value={{
          ...this.state,
          fetchKeys: this.fetchKeys,
          setKeys: this.setKeys,
        }}
      >
        {this.props.children}
      </SesKeysContext.Provider>
    );
  }
}

SesKeysProvider.propTypes = {
  children: PropTypes.element.isRequired,
};

const SesKeysConsumer = SesKeysContext.Consumer;

export { SesKeysProvider, SesKeysConsumer };
