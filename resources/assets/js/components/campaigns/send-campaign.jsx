/** @jsx React.DOM */

require('sweetalert');
var React = require('react');
var TestSend = require('./test-send-form.jsx');
var Recipients = require('./recipients.jsx');
var PreviewTemplate = require('./preview-template.jsx');
var Campaign = require('../../entities/campaign.js');
var Template = require('../../entities/template.js');

var c = new Campaign();
var t = new Template();

var SendCampaign = React.createClass({
    handleSubmit: function (e) {
        e.preventDefault();

        var lists = $('#subscribers').val();
        swal({
                title: "Are you sure?",
                text: "Do you want to send this campaign?",
                type: "info",
                showCancelButton: true,
                confirmButtonText: "Yes",
                cancelButtonText: "No",
                closeOnConfirm: false
            },
            function (isConfirm) {
                if (isConfirm) {
                    c.send(lists, this.props.data.id)
                        .done(function () {
                            swal("Sent!", "The campaign has been started!", "success");
                        }).fail(function () {
                            swal("Cancelled", "There has been an error. The campaign could not be started.", "error");
                        });

                } else {
                    swal("Cancelled", "The campaign is canceled", "error");
                }
            }.bind(this));
    },
    render: function () {
        return (
            <div className="row">
                <div className="col-lg-7 pull-right">
                    <PreviewTemplate tid={this.props.data.template_id} from={this.props.data.from_email} subject={this.props.data.subject}/>
                </div>
                <div className="row">
                    <div className="col-lg-4">
                        <TestSend cid={this.props.data.id}/>
                    </div>
                    <div className="col-lg-4">
                        <form onSubmit={this.handleSubmit}>
                            <Recipients />
                            <button type="submit" className="btn btn-default">
                                <span className="glyphicon glyphicon-envelope"></span> Send campaign
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        )
    }
});

module.exports = SendCampaign;
