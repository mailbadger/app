
import * as bootpag from 'bootpag/lib/jquery.bootpag.min.js';

import React, {Component} from 'react';
import DeleteButton from '../delete-button.jsx';
import Campaign from '../../entities/campaign.js';

const c = new Campaign();

const getSentCampaigns = (component) => {
    let data = {
        paginate: true,
        per_page: 10,
        page: 1,
        search: 'sent',
        searchFields: 'status:='
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

class ReportRow extends Component {

    constructor(props) {
        super(props);

        this.viewReport = this.viewReport.bind(this);
    }

    viewReport() {
        this.props.viewReport(this.props.data.id);
    }

    render() {
        return (
            <tr>
                <td><a href="#" onClick={this.viewReport}>{this.props.data.name}</a></td>
                <td>{this.props.data.recipients}</td>
                <td>{this.props.data.sent_at}</td>
                <td>{this.props.data.subject}</td>
                <td>{this.props.data.from_name}</td>
                <td>{this.props.data.from_email}</td>
            </tr>
        );
    }
}

export default class ReportsTable extends Component {

    constructor(props) {
        super(props);
        this.state = {
            campaigns: {
                data: []
            }
        };
    } 

    componentDidMount() {
        getSentCampaigns(this);
    }

    render() {
        let rows = (data) => {
            return <ReportRow key={data.id} data={data} viewReport={this.props.viewReport} />
        };

        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>Campaign</th>
                        <th>Recipients</th>
                        <th>Sent</th>
                        <th>Subject</th>
                        <th>From name</th>
                        <th>From email</th>
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
