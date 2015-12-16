
import _ from 'underscore';
import React, {Component} from 'react';

export default class ErrorsList extends Component {
    handleErrors(errors, i) {
        return <li key={i}>{errors[0]}</li>
    }

    render() {
        return (
            <div className="alert alert-danger alert-dismissible signin-alert" role="alert">
                <ul>{_.map(this.props.errors, this.handleErrors)}</ul>
            </div>
        );
    }
}
