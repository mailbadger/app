import React, { Component } from "react";
import PropTypes from "prop-types";

const defaultValue = {
  isOpen: false,
  message: "",
  status: "status-ok"
};

export const NotificationsContext = React.createContext(defaultValue);

class NotificationsProvider extends Component {
  constructor(props) {
    super(props);

    this.createNotification = this.createNotification.bind(this);
    this.close = this.close.bind(this);
    this.state = defaultValue;
  }

  createNotification(message, status = "status-ok") {
    this.setState({ isOpen: true, message, status });
  }

  close() {
    this.setState(defaultValue);
  }

  render() {
    return (
      <NotificationsContext.Provider
        value={{
          ...this.state,
          createNotification: this.createNotification,
          close: this.close
        }}
      >
        {this.props.children}
      </NotificationsContext.Provider>
    );
  }
}

NotificationsProvider.propTypes = {
  children: PropTypes.element.isRequired
};
const NotificationConsumer = NotificationsContext.Consumer;

export { NotificationsProvider, NotificationConsumer };
