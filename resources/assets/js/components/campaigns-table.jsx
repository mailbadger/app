/**
 * Created by filip on 22.7.15.
 */
/** @jsx React.DOM */

require('bootpag/lib/jquery.bootpag.min.js');
require('sweetalert');

var React = require('react');
var Campaign = require('../entities/campaign.js');
var c = new Campaign();

var DeleteButton = React.createClass({
    handleSubmit: function (e) {
        e.preventDefault();
        swal({
                title: "Are you sure?",
                text: "You will not be able to recover this campaign!",
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55",
                confirmButtonText: "Yes, delete it!",
                closeOnConfirm: false
            },
            function () {
                c.delete(this.props.cid)
                    .done(function () {
                        swal({
                            title: "Success",
                            text: "The campaign was successfully deleted!",
                            type: "success"
                        }, function () {
                            location.reload();
                        });
                    })
                    .fail(function () {
                        swal('Error', 'The campaign was not deleted. Try again.', 'error');
                    });
            }.bind(this));

    },
    render: function () {
        return (
            <form onSubmit={this.handleSubmit}>
                <input type="hidden" name="_method" value="DELETE"/>
                <button type="submit"><span className="delete-campaign glyphicon glyphicon-trash"></span></button>
            </form>
        );
    }
});

var CampaignRow = React.createClass({
    render: function () {
        return (
            <tr>
                <td>{this.props.data.name}</td>
                <td>{this.props.data.recipients}</td>
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
                    <DeleteButton cid={this.props.data.id}/>
                </td>
            </tr>
        );
    }
});

var CampaignsTable = React.createClass({
    getInitialState: function () {
        return {campaigns: {data: []}};
    },
    componentDidMount: function () {
        c.all(true, 10, 1).done(function (response) {
            this.setState({campaigns: response});

            $('.pagination').bootpag({
                total: response.last_page,
                page: response.current_page,
                maxVisible: 5
            }).on("page", function (event, num) {
                c.all(true, 10, num).done(function (response) {
                    this.setState({campaigns: response});
                    $('.pagination').bootpag({page: response.current_page});
                }.bind(this));
            }.bind(this));
        }.bind(this));
    },
    render: function () {
        var rows = function (data) {
            return <CampaignRow key={data.id} data={data}/>
        };
        return (
            <div>
                <table className="table table-responsive table-striped table-hover">
                    <thead>
                    <tr>
                        <th>Campaign</th>
                        <th>Recipients</th>
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
});

React.render(<CampaignsTable />, document.getElementById('campaigns'));

