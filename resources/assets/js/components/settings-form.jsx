/** @jsx React.DOM */

require('sweetalert');
var React = require('react');
var User = require('../entities/user.js');

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
            aws_key: this.refs.aws_key.getDOMNode().value,
            aws_secret: this.refs.aws_secret.getDOMNode().value,
            aws_region: this.refs.aws_region.getDOMNode().value,
        };

        var password = this.refs.password.getDOMNode().value;
        
        if(password.trim() !== '') {
            data.password = password.trim();
        }

        u.saveSettings(data).done(this.handleSuccess).fail(this.handleErrors);
    },
    render: function () {
        var errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        return (
            <div>
                <form onSubmit={this.handleSubmit}>
                    <div className="errors">{errors}</div>

                    <div className="col-sm-6 pull-left">
                        <div className="col-sm-8 form-group">
                            <label htmlFor="email">Email address</label>
                            <input type="email" className="form-control" id="email" ref="email" defaultValue={this.props.data.email} required/>
                        </div>
                        <div className="col-sm-8 form-group">
                            <label htmlFor="name">Name</label>
                            <input type="text" className="form-control" id="name" ref="name" defaultValue={this.props.data.name} required/>
                        </div>
                        <div className="col-sm-8 form-group">
                            <label htmlFor="password">Password (leave blank to not change it)</label>
                            <input type="password" className="form-control" id="password" ref="password" />
                        </div>
                    </div>
                    <div className="col-sm-6 pull-right">
                        <div className="col-sm-8 form-group">
                            <label htmlFor="aws_key">AWS Key</label>
                            <input type="text" className="form-control" id="aws_key" ref="aws_key" defaultValue={this.props.data.aws_key} required/>
                        </div>
                        <div className="col-sm-8 form-group">
                            <label htmlFor="aws_secret">AWS Secret Key</label>
                            <input type="text" className="form-control" id="aws_secret" ref="aws_secret" defaultValue={this.props.data.aws_secret} required/>
                        </div>
                        <div className="col-sm-8 form-group">
                            <label htmlFor="aws_region">AWS Region</label>
                            <select className="form-control" id="aws_region" ref="aws_region" required >
                                <option value="eu-west-1">EU (Ireland)</option> 
                                <option value="us-east-1">US East (N. Virginia)</option>
                                <option value="us-west-2">US West (Oregon)</option>
                            </select>
                        </div>
                    </div>
                    <div className="col-sm-12">
                        <button className="col-sm-2 btn btn-default">Save settings</button>
                    </div>
                </form>
            </div>
        );
    }
});

module.exports = SettingsForm;

