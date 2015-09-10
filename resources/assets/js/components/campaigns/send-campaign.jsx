/** @jsx React.DOM */

require('sweetalert');
var React = require('react');
var TestSend = require('./test-send-form.jsx');
var Campaign = require('../../entities/campaign.js');
var Template = require('../../entities/template.js');

var c = new Campaign();
var t = new Template();

var SendCampaign = React.createClass({
    render: function () {
        return (
            <div className="row">
                <TestSend cid={this.props.cid} />
            </div>
        )
    }
});

module.exports = SendCampaign;