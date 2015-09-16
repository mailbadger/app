/** @jsx React.DOM */

var _ = require('underscore');
var React = require('react');

var ErrorsList = React.createClass({
    handleErrors: function(errors, i) {
        return <li key={i}>{errors[0]}</li>
    },
    render: function () {
        return (
            <div className="alert alert-danger alert-dismissible signin-alert" role="alert">
                <ul>{_.map(this.props.errors, this.handleErrors)}</ul>
            </div>
        );
    }
});

module.exports = ErrorsList;