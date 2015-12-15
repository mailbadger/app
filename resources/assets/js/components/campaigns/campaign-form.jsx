
import sweetalert from 'sweetalert';
import React from 'react';
import Campaign from '../../entities/campaign.js';
import Template from '../../entities/template.js';
import ErrorsList from '../errors-list.jsx';

const c = new Campaign();
const t = new Template();

export default class CampaignForm extends React.Component { 

    constructor(props) {
        super(props);
        this.state = {
            hasErrors: false,
            errors: {},
            templates: []
        };

        this.handleSuccess = this.handleSuccess.bind(this);
        this.handleErrors = this.handleErrors.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSuccess() {
        this.setState({hasErrors: false, errors: []});
        sweetalert({
            title: "Success",
            text: "The campaign was successfully created!",
            type: "success"
        }, () => {
            window.location.href = url_base + '/dashboard';
        });
    }

    handleErrors(xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    }

    handleSubmit(e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});
        let data = {};

        if (!this.props.edit) {
            data = {
                name: this.refs.name.getDOMNode().value,
                subject: this.refs.subject.getDOMNode().value,
                from_name: this.refs.from_name.getDOMNode().value,
                from_email: this.refs.from_email.getDOMNode().value,
                template_id: this.refs.template.getDOMNode().value,
                status: 'draft',
                recipients: 0
            };
            c.create(data).done(this.handleSuccess).fail(this.handleErrors);
        } else {
            data = {
                name: this.refs.name.getDOMNode().value,
                subject: this.refs.subject.getDOMNode().value,
                from_name: this.refs.from_name.getDOMNode().value,
                from_email: this.refs.from_email.getDOMNode().value,
                template_id: this.refs.template.getDOMNode().value
            };
            c.update(data, this.props.data.id).done(this.handleSuccess).fail(this.handleErrors);
        }
    }

    componentDidMount() {
        t.all().done((res) => {
            this.setState({templates: res});
            if (this.props.edit) {
                $('#select-template').select2('val', this.props.data.template_id);
            }
        });

        $('#select-template').select2({placeholder: 'Select a template'});
    }

    render() {
        let errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        let backBtn = (this.props.edit) ? <a href="#" onClick={this.props.back}>Back</a> : null;
        let templates = function (t) {
            return <option value={t.id} key={t.id}>{t.name}</option>
        };

        return (
            <div>
                <form onSubmit={this.handleSubmit}>
                    <div className="errors">{errors}</div>
                    <div className="col-lg-4">
                        <div className="form-group">
                            <label htmlFor="name">Campaign name:</label>
                            <input type="text" className="form-control" ref="name" name="name" id="name"
                                   placeholder="Name" defaultValue={this.props.edit ? this.props.data.name : ''}
                                   required/>
                        </div>
                        <div className="form-group">
                            <label htmlFor="subject">Subject:</label>
                            <input type="text" className="form-control" ref="subject" name="subject" id="subject"
                                   placeholder="Subject" defaultValue={this.props.edit ? this.props.data.subject : ''}
                                   required/>
                        </div>
                        <div className="form-group">
                            <label htmlFor="from-name">From name:</label>
                            <input type="text" className="form-control" ref="from_name" name="from_name" id="from-name"
                                   placeholder="John Doe"
                                   defaultValue={this.props.edit ? this.props.data.from_name : ''} required/>
                        </div>
                        <div className="form-group">
                            <label htmlFor="from-email">From email:</label>
                            <input type="email" className="form-control" ref="from_email" name="from_email"
                                   id="from-email" defaultValue={this.props.edit ? this.props.data.from_email : ''}
                                   placeholder="example@foobar.com" required/>
                        </div>
                        <button className="col-lg-4 btn btn-default">Save campaign</button>
                    </div>
                    <div className="col-lg-6 template">
                        <label htmlFor="select-template">Select template:</label>
                        <a href={url_base + '/dashboard/new-template'}>Don't have one? Create it here.</a>
                        <select ref="template" id="select-template" style={{width: 75 + '%'}} required>
                            <option></option>
                            {this.state.templates.map(templates)}
                        </select>
                    </div>
                </form>
                <div className="col-lg-12">{backBtn}</div>
            </div>
        );
    }
}
