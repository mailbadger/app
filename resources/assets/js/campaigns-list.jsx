/** @jsx React.DOM */

var React = require('react');

var CampaignsTable = require('./components/campaigns/campaigns-table.jsx');
var CampaignForm = require('./components/campaigns/campaign-form.jsx');
var SendCampaign = require('./components/campaigns/send-campaign.jsx');
var CreateNewButton = require('./components/create-new-button.jsx');
var Campaign = require('./entities/campaign.js');

var c = new Campaign();

var Campaigns = React.createClass({
    getInitialState: function () {
        return {step: '', campaign: {}}
    },
    editCampaign: function (id) {
        c.get(id).done(function (res) {
            this.setState({step: 'edit', campaign: res});
        }.bind(this));
    },
    sendCampaign: function (id) {
        this.setState({step: 'send', cid: id});
    },
    back: function () {
        this.setState({step: ''});
    },
    render: function () {
        switch (this.state.step) {
            case 'edit':
                return <CampaignForm data={this.state.campaign} edit={true} back={this.back}/>;
            case 'send':
                return <SendCampaign cid={this.state.cid} back={this.back}/>;
            default:
                return (
                    <div>
                        <CreateNewButton url={url_base + '/dashboard/new-campaign'} text="Create new campaign"/>

                        <div className="row">
                            <CampaignsTable editCampaign={this.editCampaign} sendCampaign={this.sendCampaign}/>
                        </div>
                    </div>
                );
        }
    }

});

React.render(<Campaigns />, document.getElementById('campaigns'));