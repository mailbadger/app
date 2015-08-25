/** @jsx React.DOM */

var React = require('react');
var ErrorsList = require('../errors-list.jsx');
var List = require('../../entities/list.js');
var l = new List();

var DeleteSubscribers = React.createClass({
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
            text: "The subscribers were successfully deleted!",
            type: "success"
        }, function () {
            this.props.back();
        }.bind(this));
    },
    handleErrors: function (xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    },
    handleSubmit: function (e) {
        e.preventDefault();

        l.deleteSubscribers(this.props.listId, this.refs.subscribers.getDOMNode().files[0])
            .done(this.handleSuccess)
            .fail(this.handleErrors);
    },

    render: function () {
        var errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;

        return (
            <div>
                <div className="row">
                    <h2>Mass delete via csv/xls file</h2>

                    <p>File format:</p>
                    <ul>
                        <li>Format your CSV the same way as the example below</li>
                        <li>The number of columns in your CSV should be the same as the example below</li>
                    </ul>
                    <div className="col-lg-3">
                        <table className="table table-responsive table-hover">
                            <thead>
                            <tr>
                                <th>Email</th>
                            </tr>
                            </thead>
                            <tbody>
                            <tr>
                                <td>john@doe.com</td>
                            </tr>
                            <tr>
                                <td>jane@doe.com</td>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div className="row">
                    <form onSubmit={this.handleSubmit} id="delete-form">
                        <div className="errors">{errors}</div>
                        <div className="col-lg-3">
                        <span className="btn btn-success btn-file">
                            Browse<input type="file" ref="subscribers" name="subscribers"
                                         id="subscribers" required/>
                        </span>
                        </div>
                        <button type="submit" className="btn btn-default">Mass delete</button>
                    </form>
                    <div className="col-lg-12" style={{marginTop: '20px'}}>
                        <a href="#" onClick={this.props.back}>Back to list</a>
                    </div>
                </div>
            </div>
        );
    }
});

module.exports = DeleteSubscribers;