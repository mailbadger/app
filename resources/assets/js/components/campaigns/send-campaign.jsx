
import sweetalert from 'sweetalert';
import React, {Component} from 'react';
import Campaign from '../../entities/campaign.js';
import ErrorsList from '../errors-list.jsx';
import TestSend from './test-send-form.jsx';
import Recipients from './recipients.jsx';
import PreviewTemplate from './preview-template.jsx';

const c = new Campaign();

export default class SendCampaign extends Component {

    constructor(props) {
        super(props);

        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit(e) {
        e.preventDefault();

        let lists = $('#subscribers').val();
        sweetalert({
            title: "Are you sure?",
            text: "Do you want to send this campaign?",
            type: "info",
            showCancelButton: true,
            confirmButtonText: "Yes",
            cancelButtonText: "No",
            closeOnConfirm: false
        }, (isConfirm) => {
            if (isConfirm) {
                c.send(lists, this.props.data.id)
                .done(() => {
                    sweetalert("Sent!", "The campaign has been started!", "success");
                }).fail(() => {
                    sweetalert("Cancelled", "There has been an error. The campaign could not be started.", "error");
                });

            } else {
                sweetalert("Cancelled", "The campaign is canceled", "error");
            }
        });
    }

    render() {
        return (
            <div className="row">
                <div className="col-lg-7 pull-right">
                    <PreviewTemplate tid={this.props.data.template_id} from={this.props.data.from_email} subject={this.props.data.subject}/>
                </div>
                <div className="row">
                    <div className="col-lg-4">
                        <TestSend cid={this.props.data.id}/>
                    </div>
                    <div className="col-lg-4">
                        <form onSubmit={this.handleSubmit}>
                            <Recipients />
                            <button type="submit" className="btn btn-default">
                                <span className="glyphicon glyphicon-envelope"></span> Send campaign
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        )
    }
}
