/**
 * Created by filip on 22.7.15.
 */

import * as bootpag from 'bootpag/lib/jquery.bootpag.min.js';
import React, {Component} from 'react';
import DeleteButton from '../delete-button.jsx';
import Campaign from '../../entities/campaign.js';

const c = new Campaign();

const getAllCampaigns = (component) => {
    let data = {
        paginate: true,
        per_page: 10,
        page: 1
    };

    c.all(data).done((res) => {
        component.setState({campaigns: res});

        $('.pagination').bootpag({
            total: res.last_page,
            page: res.current_page,
            maxVisible: 5
        }).on("page", (event, num) => {
            data.page = num;
            c.all(data).done((res) => {
                component.setState({campaigns: res});
                $('.pagination').bootpag({page: res.current_page});
            });
        });
    });
}

class CampaignRow extends Component {

    constructor(props) {
        super(props);

        this.editCampaign = this.editCampaign.bind(this);
        this.sendCampaign = this.sendCampaign.bind(this);
    }

    editCampaign() {
        this.props.editCampaign(this.props.data.id);
    }
    
    sendCampaign() {
        this.props.sendCampaign(this.props.data.id);
    }
    
    render() {
        let edit = (this.props.data.status === 'draft' || this.props.data.status === 'scheduled') ?
            <span> | <a href="#" onClick={this.editCampaign}>Edit</a></span> : null;
        let campaignName = (this.props.data.status !== 'sent' || this.props.data.status === 'sending') 
                ? <span><a href="#" onClick={this.sendCampaign}>{this.props.data.name}</a>{edit}</span>
                : this.props.data.name; 

        return (
            <tr>
                <td>{campaignName}</td>
                <td>{this.props.data.subject}</td>
                <td>{this.props.data.from_name}</td>
                <td>{this.props.data.from_email}</td>
                <td>{(() => {
                    switch (this.props.data.status) {
                        case "draft":
                            return <span className="label label-default">Draft</span>;
                        case "sent":
                            return <span className="label label-success">Sent</span>;
                        case "sending":
                            return <span className="label label-info">Sending</span>;
                    }
                })()}</td>
                <td>
                    <DeleteButton success={this.props.handleDelete} resourceId={this.props.data.id} entity={c}/>
                </td>
            </tr>
        );
    }
}

export default class CampaignsTable extends Component {

    constructor(props) {
        super(props);
        this.state = {
            campaigns: {
                data: []
            }
        };

        this.handleDelete = this.handleDelete.bind(this);
    }

    componentDidMount() {
        getAllCampaigns(this);
    }

    handleDelete() {
        getAllCampaigns(this);
    }

    render() {
        let rows = (data) => {
            return <CampaignRow key={data.id} data={data} handleDelete={this.handleDelete}
                                editCampaign={this.props.editCampaign} sendCampaign={this.props.sendCampaign}/>
        }

        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>Campaign</th>
                        <th>Subject</th>
                        <th>From name</th>
                        <th>From email</th> 
                        <th>Status</th>
                        <th>Delete</th>
                    </tr>
                    </thead>
                    <tbody>
                    {this.state.campaigns.data.map(rows)}
                    </tbody>
                </table>
                <div className="col-lg-12 pagination text-center"></div>
            </div>
        );
    }
}
