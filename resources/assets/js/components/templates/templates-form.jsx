/** @jsx React.DOM */

var React = require('react');
var Template = require('../../entities/template.js');
var ErrorsList = require('../errors-list.jsx');

var t = new Template();

var TemplateForm = React.createClass({
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
            text: "The template was successfully created!",
            type: "success"
        }, function () {
            window.location.href = url_base + '/dashboard/templates';
        });
    },
    handleErrors: function (xhr) {
        this.setState({hasErrors: true, errors: xhr.responseJSON});
    },
    handleSubmit: function (e) {
        e.preventDefault();
        this.setState({hasErrors: false, errors: []});

        for (instance in CKEDITOR.instances) {
            CKEDITOR.instances[instance].updateElement();
        }

        var data = {
            name: this.refs.name.getDOMNode().value,
            content: CKEDITOR.instances.content.getData()
        };

        t.create(data).done(this.handleSuccess).fail(this.handleErrors);
    },
    componentDidMount: function () {
        CKEDITOR.replace('content');
    },
    render: function () {
        var errors = (this.state.hasErrors) ? <ErrorsList errors={this.state.errors}/> : null;
        return (
            <form onSubmit={this.handleSubmit}>
                <div className="errors">{errors}</div>
                <div className="col-lg-4">
                    <div className="form-group">
                        <label htmlFor="name">Template name:</label>
                        <input type="text" className="form-control" ref="name" name="name" id="name" placeholder="Name"
                               required/>
                    </div>
                    <button className="col-lg-4 btn btn-default">Save template</button>
                </div>
                <div className="col-lg-7">
                    <div className="form-group">
                        <label htmlFor="content">Content:</label>
                        <textarea name="content" id="content"></textarea>
                    </div>
                </div>

            </form>
        );
    }
});


React.render(<TemplateForm />, document.getElementById('new-template'));