/** @jsx React.DOM */

require('sweetalert');
var _ = require('underscore');
var React = require('react');
var Campaign = require('../../entities/campaign.js');
var c = new Campaign();

var TestSend = React.createClass({
    handleTestSend: function (e) {
        e.preventDefault();

        var emails = this.refs.emails.getDOMNode().value;
        swal({
            title: "Are you sure?",
            text: "Do you want to send this campaign to the following emails: " + emails,
            type: "info",
            showCancelButton: true,
            closeOnConfirm: false,
            showLoaderOnConfirm: true
        },
        function () {    
            c.testSend(emails.split(','), this.props.cid)
            .done(function (res) {
                swal("Sent!", "Test emails have been sent.", "success");
            }).fail(function (xhr) {
                var html = '<ul>';
                _.map(xhr.responseJSON, function(e) { 
                    html += '<li>' + e[0] + '</li>'; 
                });
                html += '</ul>';
                swal({html:true, title: "Cancelled", text: html, type: "error"});
            });
        }.bind(this));
    },
    render: function () {
        return (
            <div>
                <h3>Test send this campaign</h3>

                <form onSubmit={this.handleTestSend}>
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
