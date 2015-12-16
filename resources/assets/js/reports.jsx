
import React, {Component} from 'react';
import ReportsTable from './components/reports/reports-table.jsx';
import Report from './components/reports/report.jsx';
import Campaign from './entities/campaign.js';

const c = new Campaign();

class Reports extends Component {

    constructor(props) {
        super(props);
        this.state = {
            step: '',
            campaign: {}
        };

        this.viewReport = this.viewReport.bind(this);
        this.back = this.back.bind(this);
    }
    
    viewReport(id) {
        c.get(id).done((res) => {
            this.setState({step: 'view', campaign: res});
        }); 
    }

    back() {
        this.setState({step: ''});
    }

    render() { 
        switch (this.state.step) {
            case 'view':
                return <Report data={this.state.campaign} back={this.back}/>;
            default:
                return (
                    <div className="row">
                        <ReportsTable viewReport={this.viewReport} />
                    </div>
                ); 
        }
    }
}

React.render(<Reports />, document.getElementById('reports'));
