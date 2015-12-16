
import sweetalert from 'sweetalert';
import _ from 'underscore';
import React, {Component} from 'react';
import Campaign from '../../entities/campaign.js';

const c = new Campaign();

export default class TestSend extends Component {

    constructor(props) {
        super(props);

        this.handleTestSend = this.handleTestSend.bind(this);
    }

    handleTestSend(e) {
        e.preventDefault();

        let emails = this.refs.emails.getDOMNode().value;
        sweetalert({
            title: "Are you sure?",
            text: "Do you want to send this campaign to the following emails: " + emails,
            type: "info",
            showCancelButton: true,
            closeOnConfirm: false,
            showLoaderOnConfirm: true
        },
        function () {    
            c.testSend(emails.split(','), this.props.cid)
            .done((res) => {
                swal("Sent!", "Test emails have been sent.", "success");
            }).fail((xhr) => {
                let html = '<ul>';
                _.map(xhr.responseJSON, function(e) { 
                    html += '<li>' + e[0] + '</li>'; 
                });
                html += '</ul>';
                sweetalert({html:true, title: "Cancelled", text: html, type: "error"});
            });
        }.bind(this));
    }
    
    render() {
        return (
            <div>
                <h3>Test send this campaign</h3>

                <form onSubmit={this.handleTestSend}>
                    <div className="form-group">
                        <label htmlFor="emails">Email addresses</label>
                        <input type="text" className="form-control" ref="emails" id="emails"
                            placeholder="Emails, separated by comma"/>
                    </div>

                    <button type="submit" className="btn btn-default">
                        <span className="glyphicon glyphicon-envelope"></span> Test send this campaign
                    </button>
                </form>
            </div>
        )
    }
}
