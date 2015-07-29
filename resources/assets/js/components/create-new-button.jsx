/** @jsx React.DOM */

var React = require('react');

var CreateNewButton = React.createClass({
    render: function() {
        return (
            <div className="row">
                <div className="col-lg-4">
                    <a href={this.props.url} className="btn btn-success btn-lg">
                        <span className="glyphicon glyphicon-plus"></span> {this.props.text}
                    </a>
                </div>
            </div>
        );
    }
});

module.exports = CreateNewButton;