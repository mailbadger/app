/** @jsx React.DOM */

var React = require('react');
var Chart = require('chart.js');

var ReportInfo = React.createClass({
    render: function () {
        return (
            <div className="col-xs-4 well">
                <h4>Name: <strong>{this.props.data.name}</strong></h4>
                <h4>Subject: <strong>{this.props.data.subject}</strong></h4>
                <h4>From: <strong>{this.props.data.from_email}</strong></h4>
                <h4>Sent at: <strong>{this.props.data.sent_at}</strong></h4>
            </div>
        );
    }
});

var ReportStatistics = React.createClass({
    render: function () {
        var opensPercent = this.props.data.opens_count[0].unique_opens / this.props.data.recipients * 100;
        var bouncesPercent = this.props.data.bounces_count[0].bounces / this.props.data.recipients * 100;
        var complaintsPercent = this.props.data.complaints_count[0].complaints / this.props.data.recipients * 100;

        return (
            <div className="col-xs-6 col-xs-offset-1 well">
                <div className="row">
                    <div className="col-xs-6">
                        <h4>
                            <label className="label label-success">{opensPercent.toFixed(2)} %</label> Opened <small className="label label-default">{this.props.data.opens_count[0].unique_opens} unique / {this.props.data.opens_count[0].opens} total</small>
                        </h4>
                        <h4>
                            <label className="label label-warning">{this.props.data.recipients - this.props.data.opens_count[0].unique_opens}</label> Not opened 
                        </h4>
                    </div>
                    <div className="col-xs-6">
                        <h4>
                            <label className="label label-danger">{bouncesPercent.toFixed(2)} %</label> Bounced <small className="label label-default">{this.props.data.bounces_count[0].bounces} total</small> 
                        </h4>
                        <h4>
                            <label className="label label-inverse">{complaintsPercent.toFixed(2)} %</label> Complained <small className="label label-default">{this.props.data.complaints_count[0].complaints} total</small>  
                        </h4>
                    </div>
                </div>
            </div>
        );
    }
});

var ReportChart = React.createClass({
    componentDidMount: function () { 
        var data = {
            labels: ['Recipients', 'Bounces', 'Complaints', 'Opened', 'Unopened'],
            datasets: [
                {
                    label: 'Report dataset',
                    fillColor: "rgba(220,220,220,0.5)",
                    strokeColor: "rgba(220,220,220,0.8)",
                    highlightFill: "rgba(220,220,220,0.75)",
                    highlightStroke: "rgba(220,220,220,1)",
                    data: [
                        this.props.data.recipients,
                        this.props.data.bounces_count[0].bounces, 
                        this.props.data.complaints_count[0].complaints,
                        this.props.data.opens_count[0].unique_opens,
                        this.props.data.recipients - this.props.data.opens_count[0].unique_opens
                    ]
                }
            ]
        };

        var ctx = document.getElementById('bar-chart').getContext('2d');
        var barChart = new Chart(ctx).Bar(data);
    },
    render: function () {
        return (
            <div>
                <ReportInfo data={this.props.data}/>
                <ReportStatistics data={this.props.data}/>
                <div className="col-xs-8 col-xs-offset-2">
                    <canvas id="bar-chart" className="col-xs-12"></canvas>
                </div>
            </div>
        );
    }
});

var  Report = React.createClass({
    render: function () {
        return <ReportChart data={this.props.data}/>;
    }
});

module.exports = Report;
