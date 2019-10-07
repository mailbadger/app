import React, { Component } from "react";

const defaultValue = {
  isOpen: false,
  message: "",
  status: "status-ok"
};

export const NotificationsContext = React.createContext(defaultValue);

class NotificationsProvider extends Component {
  state = defaultValue;

  constructor(props) {
    super(props);

    this.createNotification = this.createNotification.bind(this);
    this.close = this.close.bind(this);
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

const NotificationConsumer = NotificationsContext.Consumer;

export { NotificationsProvider, NotificationConsumer };
