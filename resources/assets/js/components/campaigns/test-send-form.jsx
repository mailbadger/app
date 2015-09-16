/** @jsx React.DOM */

require('sweetalert');
var React = require('react');
var Campaign = require('../../entities/campaign.js');
var c = new Campaign();

var TestSend = React.createClass({
    handleSubmit: function () {
        var emails = this.refs.emails.getDOMNode().value;
        swal({
                title: "Are you sure?",
                text: "Do you want to send this campaign to the following emails: " + emails,
                type: "info",
                showCancelButton: true,
                confirmButtonText: "Yes",
                cancelButtonText: "No",
                closeOnConfirm: false
            },
            function (isConfirm) {
                if (isConfirm) {
                    c.testSend(emails.split(','), this.props.cid)
                        .done(function () {
                            swal("Sent!", "Test emails have been sent.", "success");
                        }).fail(function () {
                            swal("Cancelled", "Test emails could not be sent, check the input if they are in a correct format. Use commas for separation.", "error");
                        });

                } else {
                    swal("Cancelled", "Test emails have been canceled", "error");
                }
            }.bind(this));
    },
    render: function () {
        return (
            <div>
                <h3>Test send this campaign</h3>

                <form onSubmit={this.handleSubmit}>
                    <div className="form-group">
                        <label htmlFor="emails">Email addresses</label>
                        <input type="text" className="form-control" ref="emails" id="emails"
                               placeholder="Emails, separated by comma"/>
                    </div>

                    <button type="submit" className="btn btn-default">
                        <span className="glyphicon glyphicon-envelope"></span> Test send this campaign
                    </button>
                </form>
            </div>
        )
    }
});

module.exports = TestSend;