/** @jsx React.DOM */

var React = require('react');

var TemplateForm = require('./components/templates/template-form.jsx');
var TemplatesTable = require('./components/templates/templates-table.jsx');
var CreateNewButton = require('./components/create-new-button.jsx');
var Template = require('./entities/template.js');

var t = new Template();

var Templates = React.createClass({
    getInitialState: function () {
        return {step: '', template: {}}
    },
    editTemplate: function (id) {
        t.get(id).done(function(res) {
            this.setState({step: 'edit', template: res});
        }.bind(this));
    },
    back: function () {
        this.setState({step: ''});
    },
    render: function () {
        switch (this.state.step) {
            case 'edit':
                return <TemplateForm data={this.state.template} edit={true} back={this.back}/>;
            default:
                return (
                    <div>
                        <CreateNewButton url={url_base + '/dashboard/new-template'} text="Create new template" />
                        <div className="row">
                            <TemplatesTable editTemplate={this.editTemplate} />
                        </div>
                    </div>
                );
        }
    }
});

React.render(<Templates />, document.getElementById('templates'));