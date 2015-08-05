/** @jsx React.DOM */

require('sweetalert');
var React = require('react');
var List = require('../../entities/list.js');
var ErrorsList = require('../errors-list.jsx');

var l = new List();

var ListForm = React.createClass({
    getInitialState: function () {
        return {
            hasErrors: false,
            errors: {}
        };
    },
    handleSuccess: function () {
        this.setState({hasErrors: false, errors: []});
        swal({
            title: "Success",
            text: "The list was successfully created!",
            type: "success"
        }, function () {
            window.location.href = url_base + '/dashboard/subscribers';
        });
    },
    handleErrors: function (xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    },
    handleSubmit: function (e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});
        var data;

        if(!this.props.edit) {
            data = {
                name: this.refs.name.getDOMNode().value,
                total_subscribers: 0
            };
            l.create(data).done(this.handleSuccess).fail(this.handleErrors);
        } else {
            data = {
                name: this.refs.name.getDOMNode().value
            };
            l.update(data, this.props.data.id).done(this.handleSuccess).fail(this.handleErrors);
        }
    },
    render: function () {
        var errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        var backBtn = (this.props.edit) ? <a href="#" onClick={this.props.back}>Back</a> : null;
        return (
            <form onSubmit={this.handleSubmit}>
                <div className="errors">{errors}</div>
                <div className="col-lg-4">
                    <div className="form-group">
                        <label htmlFor="name">Template name:</label>
                        <input type="text" className="form-control" ref="name" name="name" id="name" placeholder="Name"
                               defaultValue={this.props.edit ? this.props.data.name : ''} required/>
                    </div>
                    <button className="col-lg-4 btn btn-default">Save list</button>
                </div>
                <div className="col-lg-12">{backBtn}</div>
            </form>
        );
    }
});

module.exports = ListForm;