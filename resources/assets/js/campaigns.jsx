
import React, {Component} from 'react';

import Campaign from './entities/campaign.js';
import CampaignsTable from './components/campaigns/campaigns-table.jsx';
import CampaignForm from './components/campaigns/campaign-form.jsx';
import SendCampaign from './components/campaigns/send-campaign.jsx';
import CreateNewButton from './components/create-new-button.jsx';

const c = new Campaign();

class Campaigns extends Component {

    constructor(props) {
        super(props);
        this.state = {
            step: '',
            campaign: {}
        };

        this.editCampaign = this.editCampaign.bind(this);
        this.sendCampaign = this.sendCampaign.bind(this);
        this.back = this.back.bind(this);
    }

    editCampaign(id) {
        c.get(id).done((res) => {
            this.setState({step: 'edit', campaign: res});
        });
    }

    sendCampaign(id) {
        c.get(id).done((res) => {
            this.setState({step: 'send', campaign: res});
        });
    }

    back() {
        this.setState({step: ''});
    }

    render() {
        switch (this.state.step) {
            case 'edit':
                return <CampaignForm data={this.state.campaign} edit={true} back={this.back}/>;
            case 'send':
                return <SendCampaign data={this.state.campaign} back={this.back}/>;
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
}

React.render(<Campaigns />, document.getElementById('campaigns'));
