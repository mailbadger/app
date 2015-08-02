/** @jsx React.DOM */

var React = require('react');

var CampaignForm = require('./components/campaigns/campaign-form.jsx');

React.render(<CampaignForm edit={false} />, document.getElementById('new-campaign'));