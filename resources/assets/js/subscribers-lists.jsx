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
        return {
            content: <div>
                <CreateNewButton url={url_base + '/dashboard/new-subscribers'} text="Create new list"/>

                <div className="row">
                    <ListsTable showList={this.showList} editList={this.editList}/>
                </div>
            </div>
        }
    },
    showList: function (id) {
        l.get(id).done(function (res) {
            this.setState({
                content: <SubscribersList list={res} editList={this.editList} customFields={this.customFields}
                                          back={this.back}/>
            });
        }.bind(this));
    },
    editList: function (id) {
        l.get(id).done(function (res) {
            this.setState({content: <ListForm data={res} edit={true} back={this.back}/>});
        }.bind(this));
    },
    customFields: function (id) {
        this.setState({content: <CustomFields listId={id} back={this.back}/>});
    },
    back: function () {
        this.setState({
            content: <div>
                <CreateNewButton url={url_base + '/dashboard/new-subscribers'} text="Create new list"/>

                <div className="row">
                    <ListsTable showList={this.showList} editList={this.editList}/>
                </div>
            </div>
        });
    },
    render: function () {
        return this.state.content;
    }
});

React.render(<Lists />, document.getElementById('sub-lists'));