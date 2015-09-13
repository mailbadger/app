/** @jsx React.DOM */

require('sweetalert');
var React = require('react');
var TestSend = require('./test-send-form.jsx');
var PreviewTemplate = require('./preview-template.jsx');
var Campaign = require('../../entities/campaign.js');
var Template = require('../../entities/template.js');

var c = new Campaign();
var t = new Template();

var SendCampaign = React.createClass({
    render: function () {
        return (
            <div className="row">
                <div className="col-lg-4">
                    <TestSend cid={this.props.data.id}/>
                </div>
                <div className="col-lg-7 col-lg-offset-1">
                    <PreviewTemplate tid={this.props.data.template_id} from={this.props.data.from_email} subject={this.props.data.subject}/>
                </div>
            </div>
        )
    }
});

module.exports = SendCampaign;