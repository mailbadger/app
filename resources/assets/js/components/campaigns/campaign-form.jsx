/** @jsx React.DOM */
require('sweetalert');

var React = require('react');
var Campaign = require('../../entities/campaign.js');
var ErrorsList = require('../errors-list.jsx');

var c = new Campaign();

var CampaignForm = React.createClass({
    getInitialState: function () {
        return {
            hasErrors: false,
            errors: {}
        };
    },
    handleSuccess: function () {
        this.setState({hasErrors: false, errors: []});
        swal({
            title: "Success",
            text: "The campaign was successfully created!",
            type: "success"
        }, function () {
            window.location.href = url_base + '/dashboard';
        });
    },
    handleErrors: function (xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    },
    handleSubmit: function (e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});
        var data = {
            name: this.refs.name.getDOMNode().value,
            subject: this.refs.subject.getDOMNode().value,
            from_name: this.refs.from_name.getDOMNode().value,
            from_email: this.refs.from_email.getDOMNode().value,
            template_id: this.refs.template.getDOMNode().value,
            status: 'draft',
            recipients: 0
        };
        c.create(data).done(this.handleSuccess).fail(this.handleErrors);
    },
    componentDidMount: function () {
        $('#select-template').select2({
            placeholder: 'Select a template',
            minimumResultsForSearch: Infinity,
            ajax: {
                url: url_base + '/api/templates',
                dataType: 'json',
                processResults: function (templates) {
                    var results = templates.map(function (t) {
                        return {'text': t.name, 'id': t.id};
                    });
                    return {
                        results: results
                    };
                },
                cache: true
            }
        });
    },
    render: function () {
        var errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors} /> : null;
        return (
            <form onSubmit={this.handleSubmit}>
                <div className="errors">{errors}</div>
                <div className="col-lg-4">
                    <div className="form-group">
                        <label htmlFor="name">Campaign name:</label>
                        <input type="text" className="form-control" ref="name" name="name" id="name" placeholder="Name"
                               required/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="subject">Subject:</label>
                        <input type="text" className="form-control" ref="subject" name="subject" id="subject"
                               placeholder="Subject" required/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="from-name">From name:</label>
                        <input type="text" className="form-control" ref="from_name" name="from_name" id="from-name"
                               placeholder="John Doe" required/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="from-email">From email:</label>
                        <input type="email" className="form-control" ref="from_email" name="from_email" id="from-email"
                               placeholder="example@foobar.com" required/>
                    </div>
                    <button className="col-lg-4 btn btn-default">Save campaign</button>
                </div>
                <div className="col-lg-6 template">
                    <label htmlFor="select-template">Select template:</label>
                    <a href={url_base + '/dashboard/new-template'}>Don't have one? Create it here.</a>
                    <select ref="template" id="select-template" style={{width: 75 + '%'}} required>
                        <option></option>
                    </select>
                </div>
            </form>
        );
    }
});


React.render(<CampaignForm />, document.getElementById('new-campaign'));