/** @jsx React.DOM */

var React = require('react');

var SubscribersList = require('./components/sub-lists/list.jsx');
var ListForm = require('./components/sub-lists/list-form.jsx');
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
    showList: function(id) {
        l.get(id).then(function(res) {
            var list = res;
            l.getSubscribers(id, true, 10).done(function(res) {
                this.setState({content: <SubscribersList list={list} subscribers={res} editList={this.editList} back={this.back} />});
            }.bind(this));
        }.bind(this));
    },
    editList: function (id) {
        l.get(id).done(function (res) {
            this.setState({content: <ListForm data={res} edit={true} back={this.back}/>});
        }.bind(this));
    },
    back: function () {
        this.setState({
            content: <div>
                <CreateNewButton url={url_base + '/dashboard/new-subscribers'} text="Create new list"/>

                <div className="row">
                    <ListsTable showList={this.showList} editList={this.editList} />
                </div>
            </div>
        });
    },
    render: function () {
        return this.state.content;
    }
});

React.render(<Lists />, document.getElementById('sub-lists'));