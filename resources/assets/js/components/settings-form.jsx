/** @jsx React.DOM */

require('sweetalert');
var React = require('react');
var User = require('../../entities/user.js');

var u = new User();

var SettingsForm = React.createClass({
    getInitialState: function () {
        return {
            hasErrors: false,
            errors: {},
        };
    },
    handleSuccess: function () {
        this.setState({hasErrors: false, errors: []});
        swal({
            title: "Success",
            text: "The settings have been saved",
            type: "success"
        });
    },
    handleErrors: function (xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    },
    handleSubmit: function (e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});
        var data = {
            name: this.refs.name.getDOMNode().value,
            email: this.refs.email.getDOMNode().value,
            password: this.refs.password.getDOMNode().value, 
            aws_key: this.refs.aws_key.getDOMNode().value,
            aws_secret: this.refs.aws_secret.getDOMNode().value,
            aws_region: this.refs.aws_region.getDOMNode().value,
        };

        u.saveSettings(data).done(this.handleSuccess()).fail(this.handleErrors());
    },
    render: function () {
        var errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        return (
            <div>
                <form onSubmit={this.handleSubmit}>
                    <div className="errors">{errors}</div>
                    <div class="col-sm-6 form-group">
                        <label for="email">Email address</label>
                        <input type="email" class="form-control" id="email" ref="email" defaultValue={this.props.data.email} required/>
                    </div>
                    <div class="col-sm-6 form-group">
                        <label for="name">Name</label>
                        <input type="text" class="form-control" id="name" ref="name" defaultValue={this.props.data.name} required/>
                    </div>
                    <div class="col-sm-6 form-group">
                        <label for="password">Password (leave blank to not change it)</label>
                        <input type="password" class="form-control" id="password" ref="password" />
                    </div>

                    <button className="col-lg-4 btn btn-default">Save settings</button>
                </form>
            </div>
        );
    }
});

module.exports = SettingsForm;

