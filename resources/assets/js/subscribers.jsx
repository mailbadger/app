/** @jsx React.DOM */

var React = require('react');

var SubscribersList = require('./components/sub-lists/list.jsx');
var ListForm = require('./components/sub-lists/list-form.jsx');
var CustomFields = require('./components/sub-lists/custom-fields.jsx');
var ListsTable = require('./components/sub-lists/lists-table.jsx');
var CreateNewButton = require('./components/create-new-button.jsx');
var List = require('./entities/list.js');

var l = new List();

var Lists = React.createClass({
    getInitialState: function () {
        return {step: '', list: {}};
    },
    showList: function (id) {
        l.get(id).done(function (res) {
            this.setState({step: 'show', list: res});
        }.bind(this));
    },
    editList: function (id) {
        l.get(id).done(function (res) {
            this.setState({step: 'edit', list: res});
        }.bind(this));
    },
    customFields: function (id) {
        this.setState({step: 'custom-fields'});
    },
    back: function () {
        this.setState({step: ''});
    },
    render: function () {
        switch (this.state.step) {
            case 'show':
                return <SubscribersList list={this.state.list} editList={this.editList} customFields={this.customFields}
                    back={this.back}/>;
            case 'edit':
                return <ListForm data={this.state.list} edit={true} back={this.back}/>;
            case 'custom-fields':
                return <CustomFields listId={this.state.list.id} back={this.back}/>;
            default:
                return (
                    <div>
                        <CreateNewButton url={url_base + '/dashboard/new-subscribers'} text="Create new list"/>

                        <div className="row">
                            <ListsTable showList={this.showList} editList={this.editList}/>
                        </div>
                    </div>
                );     
        }
    }
});

React.render(<Lists />, document.getElementById('sub-lists'));
